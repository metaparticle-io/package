(function () {
    var shell = require('shelljs');
    var fs = require('fs');
    
    var ports = (options) => {
        var portArr = [];
        if (options && options.ports) {
            for (var ix = 0; ix < options.ports.length; ix++) {
                portArr.push({"number": options.ports[ix]});
            }
        }
        return portArr;
    };

    module.exports.run = (img, name, options) => {
        var service = {
            "name": name,
            "guid": 1234567, 
            "services": [ 
                {
                    "name": name,
                    "replicas": options.replicas,
                    "shardSpec": options.shardSpec,
                    "containers": [
                        { "image": img }
                    ],
                    "ports": ports(options)
                }
            ],
            "serve": {
                "name": name,
            }
        };
        if (options && options.public) {
            service["serve"]["public"] = true;
        }

        if (!fs.existsSync(".metaparticle")) {
            fs.mkdirSync(".metaparticle");
        }
        var serviceString = JSON.stringify(service, null, '\t');
        fs.writeFileSync(".metaparticle/service.json", serviceString);

        shell.exec('mp-compiler -f .metaparticle/service.json');
        shell.exec('mp-compiler -f .metaparticle/service.json --deploy=false --attach=true');
    };

    module.exports.cancel = (name) => {
        shell.exec('mp-compiler -f .metaparticle/service.json --delete');
    };
})();
