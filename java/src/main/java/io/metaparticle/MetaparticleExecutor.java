package io.metaparticle;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.microsoft.rest.serializer.JacksonAdapter;
import io.metaparticle.annotations.Runtime;
import io.metaparticle.models.Container;
import io.metaparticle.models.EnvVar;
import io.metaparticle.models.Service;
import io.metaparticle.models.ServicePort;
import io.metaparticle.models.ServeSpecification;
import io.metaparticle.models.ServiceSpecification;
import java.io.IOException;
import io.metaparticle.models.ShardSpecification;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

import static io.metaparticle.Util.addAll;
import static io.metaparticle.Util.handleErrorExec;

public class MetaparticleExecutor implements Executor {
    Path specPath;

    public boolean cancel(String id) {
        return handleErrorExec(new String[] { "mp-compiler", "-f", specPath.toString(), "--delete" }, null, null);
    }

    public boolean logs(String id, OutputStream stdout, OutputStream stderr) {
        return handleErrorExec(new String[] {"mp-compiler", "-f", specPath.toString(), "--deploy=false", "--attach=true"}, stdout, stderr);
    }

    private Service makeShardedService(String image, String name, Runtime config) {
        ServeSpecification serveSpec = new ServeSpecification()
            .withName(name)
            .withPublicProperty(true);
    
        List<EnvVar> envList = new ArrayList<>();
        envList.add(new EnvVar().withName("METAPARTICLE_IN_CONTAINER").withValue("true"));
    
        List<Container> containerList = new ArrayList<>();
        containerList.add(new Container().withImage(image).withEnv(envList));
    
        List<ServicePort> portList = new ArrayList<>();
        portList.add(new ServicePort().withNumber(config.ports()[0]).withProtocol("TCP"));
    
        ShardSpecification shardSpec = new ShardSpecification()
            .withShards(config.shards())
            .withUrlPattern(config.urlShardPattern());

        List<ServiceSpecification> serviceList = new ArrayList<>();
        serviceList.add(new ServiceSpecification()
            .withName(name)
            .withShardSpec(shardSpec)
            .withContainers(containerList)
            .withPorts(portList));

        Service s = new Service()
            .withName(name)
            .withGuid(1234567)
            .withServices(serviceList)
            .withServe(serveSpec);

        return s;
    }

    private Service makeReplicatedService(String image, String name, Runtime config) {
        ServeSpecification serveSpec = new ServeSpecification()
        .withName(name)
        .withPublicProperty(true);
    
        List<EnvVar> envList = new ArrayList<>();
        envList.add(new EnvVar().withName("METAPARTICLE_IN_CONTAINER").withValue("true"));
    
        List<Container> containerList = new ArrayList<>();
        containerList.add(new Container().withImage(image).withEnv(envList));
    
        List<ServicePort> portList = new ArrayList<>();
        portList.add(new ServicePort().withNumber(config.ports()[0]).withProtocol("TCP"));
    
        List<ServiceSpecification> serviceList = new ArrayList<>();
        serviceList.add(new ServiceSpecification()
        .withName(name)
        .withReplicas(config.replicas())
        .withContainers(containerList)
        .withPorts(portList));

        Service s = new Service()
        .withName(name)
        .withGuid(1234567)
        .withServices(serviceList)
        .withServe(serveSpec);

        return s;
    }

    public boolean run(String image, String name, Runtime config, OutputStream stdout, OutputStream stderr) {
        Service s;
        if (config.shards() > 0) {
            s = makeShardedService(image, name, config);
        } else {
            s = makeReplicatedService(image, name, config);
        }

        try {
            Path cwd = Paths.get(".");
            Path dir = cwd.resolve(".metaparticle");
            Files.createDirectories(dir);
            specPath = dir.resolve("spec.json");
            Gson gson = new GsonBuilder().setPrettyPrinting().create();
            String json = gson.toJson(s);
            Files.write(specPath, json.getBytes());
            return handleErrorExec(new String[] { "mp-compiler", "-f", specPath.toString() }, stdout, stderr);
        } catch (IOException ex) {
            throw new IllegalStateException(ex);
        }
    }
}
