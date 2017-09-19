package io.metaparticle;

import java.io.IOException;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Path;

import io.metaparticle.annotations.Runtime;

import static io.metaparticle.Util.addAll;
import static io.metaparticle.Util.handleErrorExec;

public class MetaparticleExecutor implements Executor {
    public boolean cancel(String id) {
        handleErrorExec(new String[] {"kubectl", "delete", "deployments", id}, null, null);
        return handleErrorExec(new String[] {"kubectl", "delete", "services", id}, null, null);
    }

    public boolean logs(String id, OutputStream stdout, OutputStream stderr) {
        return handleErrorExec(new String[] {"ktail", "-l", "app=%s", "--template", "\"{{.Message}}\\"}, stdout, stderr);
    }

    public boolean run(String image, String name, Runtime config, OutputStream stdout, OutputStream stderr) {
            String spec = "" /*{" +
    "\"name\": \"name\"," +
    "\"guid\": 1234567," +
    "\"services\": ["
    "    {{
    "           \"name\": \"{0}\",
    "        \"replicas\": {1},
            \"containers\": [
                {{
                    \"image\": \"{2}\",
                    \"env\": [{{
                        \"name\": \"METAPARTICLE_IN_CONTAINER\",
                        \"value\": \"true\"
                    }}]
                }}
            ],
            \"ports\": [{{
                \"number\": {3}
            }}]
        }}
    ],
    \"serve\": {{
        \"name\": \"{0}\",
        \"public\": true
    }}
}}";*/
            try {
                Path specPath = Files.createTempFile("spec", "json");
                Files.write(specPath, String.format(spec, name, config.replicas(), image, config.ports()[0]).getBytes());

                return handleErrorExec(new String[] {"/home/bburns/gopath/bin/compiler", "-f", specPath.toString()}, stdout, stderr);
            } catch (IOException ex) {
                throw new IllegalStateException(ex);
            }
        }
    }