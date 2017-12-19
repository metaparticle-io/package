(function() {
    var fs = require('fs');
    var path = require('path');

    inDockerContainer = () => {
        switch (process.env.METAPARTICLE_IN_CONTAINER) {
            case 'true':
            case '1':
                return true;
            case 'false':
            case '0':
                return false;
        }
        try {
            var info = fs.readFileSync("/proc/1/cgroup");
        } catch (err) {
            return false;
        }
        // This is a little approximate...
        if (info.indexOf("docker") != -1) {
            return true;
        }
        if (info.indexOf("kubepods") != -1) {
            return true;
        }
        return false;
    };

    writeDockerfile = (options) => {
        if (options.dockerfile) {
            fs.copyFileSync(src, 'Dockerfile');
            return;
        }
        var name = options.name;
        var dockerfile = `FROM node:6-alpine
        
        COPY ./ /${name}/
        RUN npm --prefix /${name}/ install
        
        CMD npm --prefix /${name}/ start
        `;

        fs.writeFileSync('Dockerfile', dockerfile);
    };

    selectBuilder = (buildSpec) => {
        switch (buildSpec) {
            case 'docker':
                return require('./docker-builder');
            default:
                throw `Unknown builder: ${buildSpec}`;
        }
    }

    selectRunner = (execSpec) => {
        switch (execSpec) {
            case 'docker':
                return require('./docker-runner');
            case 'aci':
                return require('./aci-runner');
            case 'metaparticle':
                return require('./metaparticle-runner.js');
            default:
                throw `Unknown runner: ${execSpec}`;
        }
    }

    module.exports.containerize = (arg1, arg2) => {
        var options = arg1;
        var fn = arg2;

        if (typeof arg1 === "function") {
            fn = arg1;
            options = null;
        }
        if (inDockerContainer()) {
            fn();
        } else {
            var dir = process.cwd();
            var pkgJson = JSON.parse(fs.readFileSync(path.join(dir, 'package.json')));
            var name = pkgJson.name;
            var img = name;

            var builder = selectBuilder((options && options.builder) ? options.builder : 'docker');
            var runner = selectRunner((options && options.runner) ? options.runner : 'docker');

            process.on('SIGINT', function() {
                runner.cancel(name);
                process.exit();
            });
            if (options && options.repository) {
                img = options.repository + '/' + img;
            }
            writeDockerfile(options);
            builder.build(img);
            if (options && options.publish) {
                builder.publish(img);
            }
            runner.run(img, name, options);
        }
    };
})();