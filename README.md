# kmigrate

```shell
kmigrate is declarative patcher fo Kubernetes resources

Usage:
  kmigrate [command]

Available Commands:
  create      create patch from provided mergefile
  help        Help about any command
  merge       Merge resources defined in patch file
  version     Print the version number of kmigrate

Flags:
  -h, --help      help for kmigrate
      --verbose   Show debugging information

Use "kmigrate [command] --help" for more information about a command.

```

Look at `example/cronjobs.yaml` for example how migration file looks like. There is also interactive `kmigrate create` which guides you trought creation of the migration.


Example of migration:

<pre>
Generating patch

--- /apis/batch/v1beta1/namespaces/ff42a2e5-0dd3-46df-8334-20ef21ffdb64/cronjobs/web-backup
+++ Patched
@@ -63,6 +63,8 @@
             - mountPath: /var/www
               name: data-web
               subPath: www
+            - mountPath: /var/www/cache
+              name: cache
             - mountPath: /var/www/auth
               name: auth
               readOnly: true
@@ -75,6 +77,8 @@
           - name: data-web
             persistentVolumeClaim:
               claimName: store-storage
+          - emptyDir: {}
+            name: cache
           - name: auth
             secret:
               defaultMode: 420

Press y to patch, n to skip!

y

Patching
--- Original
+++ /apis/batch/v1beta1/namespaces/ff42a2e5-0dd3-46df-8334-20ef21ffdb64/cronjobs/web-backup
@@ -1,3 +1,5 @@
+apiVersion: batch/v1beta1
+kind: CronJob
 metadata:
   annotations:
     app.woocart.com/domain: astra.mywoocart.com
@@ -9,7 +11,7 @@
     woocart: web-backup
   name: web-backup
   namespace: ff42a2e5-0dd3-46df-8334-20ef21ffdb64
-  resourceVersion: &quot;51153037&quot;
+  resourceVersion: &quot;51442926&quot;
   selfLink: /apis/batch/v1beta1/namespaces/ff42a2e5-0dd3-46df-8334-20ef21ffdb64/cronjobs/web-backup
   uid: 804e13ec-0f45-11e9-8113-42010a9c00f5
 spec:
@@ -63,6 +65,8 @@
             - mountPath: /var/www
               name: data-web
               subPath: www
+            - mountPath: /var/www/cache
+              name: cache
             - mountPath: /var/www/auth
               name: auth
               readOnly: true
@@ -75,6 +79,8 @@
           - name: data-web
             persistentVolumeClaim:
               claimName: store-storage
+          - emptyDir: {}
+            name: cache
           - name: auth
             secret:
               defaultMode: 420
@@ -85,4 +91,6 @@
   schedule: &apos;@daily&apos;
   successfulJobsHistoryLimit: 3
   suspend: false
+status:
+  lastScheduleTime: &quot;2019-01-27T00:00:00Z&quot;

</pre>
