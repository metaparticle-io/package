(function() {
    var shell = require('shelljs');
    var fs = require('fs');

    var ports = (options) => {
        var portArr = [];
        if (options && options.ports) {
            for (var ix = 0; ix < options.ports.length; ix++) {
                portArr.push({ "number": options.ports[ix] });
            }
        }
        return portArr;
    };

    module.exports.run = (img, name, options) => {
        var service = {
            "name": name,
            "guid": 1234567,
        };

        if (options && (options.replicas || options.shardSpec)) {
            service.services = [{
                "name": name,
                "replicas": options.replicas,
                "shardSpec": options.shardSpec,
                "containers": [
                    { "image": img }
                ],
                "ports": ports(options)
            }];
            service.serve = {
                "name": name,
            };
            if (options && options.public) {
                service.serve.public = true;
            }
        }
        if (options && options.jobSpec) {
            service.jobs = [{
                "name": name,
                "replicas": options.jobSpec.count,
                "containers": [
                    { "image": img }
                ]
            }];
        }
        if (!fs.existsSync(".metaparticle")) {
            fs.mkdirSync(".metaparticle");
        }
        var serviceString = JSON.stringify(service, null, '\t');
        fs.writeFileSync(".metaparticle/service.json", serviceString);

        // TODO: these should all be async so that the ctrl-c handling above works properly...
        shell.exec('mp-compiler -f .metaparticle/service.json');
        shell.exec('mp-compiler -f .metaparticle/service.json --deploy=false --attach=true');
    };

    module.exports.cancel = () => {
        // TODO: this should be async.
        var res = shell.exec('mp-compiler -f .metaparticle/service.json --delete');
        if (res.code != 0) {
            console.log('delete failed!');
            console.log(res.stdout);
            console.log(res.stderr);
        }
    };
})();