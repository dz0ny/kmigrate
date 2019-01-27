package migrate

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kmigrate/logger"
	"os"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/ghodss/yaml"
	createjsonpatch "github.com/mattbaird/jsonpatch"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/sergi/go-diff/diffmatchpatch"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

var log = logger.New("kmigrate.migrate")

type Migrate struct {
	config *Config
	dryRun bool
	client *rest.RESTClient
}

func NewMigrate(filename string, dryRun bool, restClient *rest.RESTClient) *Migrate {
	config, err := NewConfigFromFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return &Migrate{config, dryRun, restClient}
}

func (m *Migrate) Run() {
	list := Response{}
	req := m.client.Get().
		Prefix(m.config.GetAPIPath()).
		Resource(m.config.GetResource()).
		Timeout(10 * time.Second)

	selector := m.config.GetLabelSelector()
	if selector != "" {
		req.Param("labelSelector", selector)
	}

	log.Infof("Requesting %s", req.URL())

	resp, err := req.Do().Raw()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(resp, &list)
	if err != nil {
		log.Fatal(err)
	}

	patch, err := m.config.GetPatch()
	if err != nil {
		log.Fatal(err)
	}

	decodedPatch, err := jsonpatch.DecodePatch(patch)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range list.Items {
		resource := item
		originalJSONSpec, err := json.Marshal(resource.Spec)
		if err != nil {
			log.Fatal(err)
		}

		modifiedJSONSpec, err := decodedPatch.Apply(originalJSONSpec)
		if err != nil {
			log.Fatal(err)
		}
		var spec interface{}
		err = json.Unmarshal(modifiedJSONSpec, &spec)

		resource.Spec = spec

		originalJSON, err := json.Marshal(item)
		if err != nil {
			log.Fatal(err)
		}
		modifiedJSON, err := json.Marshal(resource)
		if err != nil {
			log.Fatal(err)
		}
		modifiedYAML, err := yaml.JSONToYAML(modifiedJSON)
		if err != nil {
			log.Fatal(err)
		}
		originalYAML, err := yaml.JSONToYAML(originalJSON)
		if err != nil {
			log.Fatal(err)
		}

		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(string(originalYAML)),
			B:        difflib.SplitLines(string(modifiedYAML)),
			FromFile: item.SelfLink,
			ToFile:   "Patched",
			Context:  3,
		}
		text, _ := difflib.GetUnifiedDiffString(diff)
		fmt.Println(text)

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Press y to patch, n to skip!")
			if r, len, _ := reader.ReadRune(); len > 0 {
				if r == 'y' {
					fmt.Println("Patching")
					res := m.client.Patch(types.StrategicMergePatchType).Body(modifiedJSON).RequestURI(item.SelfLink).Timeout(10 * time.Second).Do()
					ress, err := res.Raw()
					if err != nil {
						log.Fatal(err)
					}
					newYAML, err := yaml.JSONToYAML(ress)
					if err != nil {
						log.Fatal(err)
					}
					diff := difflib.UnifiedDiff{
						A:        difflib.SplitLines(string(originalYAML)),
						B:        difflib.SplitLines(string(newYAML)),
						FromFile: "Original",
						ToFile:   item.SelfLink,
						Context:  3,
					}
					text, _ := difflib.GetUnifiedDiffString(diff)
					fmt.Println(text)
				} else {
					fmt.Println("Skipping")
				}
				break
			}
		}
	}
}

func (m *Migrate) Create() {
	list := Response{}
	req := m.client.Get().
		Prefix(m.config.GetAPIPath()).
		Resource(m.config.GetResource()).
		Timeout(10 * time.Second)

	selector := m.config.GetLabelSelector()
	if selector != "" {
		req.Param("labelSelector", selector)
	}

	log.Infof("Requesting %s", req.URL())

	resp, err := req.Do().Raw()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(resp, &list)
	if err != nil {
		log.Fatal(err)
	}

	f, err := ioutil.TempFile("", "*.yaml")
	defer os.Remove(f.Name())

	originalYAML, err := yaml.Marshal(list.Items[0].Spec)
	if err != nil {
		log.Fatal(err)
	}
	f.Write(originalYAML)
	f.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Open file %s in editor to create patch.\n", f.Name())
	for {
		fmt.Print("Press enter when done!")
		if _, len, _ := reader.ReadRune(); len > 0 {
			break
		}
	}

	originalJSON, err := json.Marshal(list.Items[0].Spec)
	if err != nil {
		log.Fatal(err)
	}
	modifiedYAML, _ := ioutil.ReadFile(f.Name())
	modifiedJSON, err := yaml.YAMLToJSON(modifiedYAML)
	if err != nil {
		log.Fatalf("Could not parse patch file, make sure it's valid YAML. %v", err)
	}

	patch, err := createjsonpatch.CreatePatch(originalJSON, modifiedJSON)
	if err != nil {
		log.Fatal(err)
	}
	jsonPatch, err := json.Marshal(patch)
	if err != nil {
		log.Fatal(err)
	}
	yamlPatchRaw, err := yaml.JSONToYAML(jsonPatch)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generated patch:")
	fmt.Println(string(jsonPatch))
	fmt.Println(string(yamlPatchRaw))

	if err = m.config.UpdatePatch(patch); err != nil {
		log.Fatal(err)
	}
	generateDiff(originalJSON, jsonPatch)
}

func generateDiff(originalJSON, patch []byte) {

	decodedPatch, err := jsonpatch.DecodePatch(patch)
	if err != nil {
		log.Fatal(err)
	}

	modified, err := decodedPatch.Apply(originalJSON)
	if err != nil {
		log.Fatal(err)
	}

	originalYAML, err := yaml.JSONToYAML(originalJSON)
	if err != nil {
		log.Fatal(err)
	}

	modifiedYAML, err := yaml.JSONToYAML(modified)
	if err != nil {
		log.Fatal(err)
	}
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(string(originalYAML), string(modifiedYAML), true)

	fmt.Println(dmp.DiffPrettyText(diffs))
}
