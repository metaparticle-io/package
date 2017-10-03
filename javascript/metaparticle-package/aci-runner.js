(function () {
    var shell = require('shelljs');
    var rg = 'test';
    
    portString = (options) => {
        if (!options || !options.ports) {
            return '';
        }
        return `--port=${options.ports[0]}`;
    };

    module.exports.run = (img, name, options) => {
        var ports = portString(options);
        var public = '';
        if (options && options.public) {
            public = '--ip-address Public';
        }
        shell.exec(`az container create ${public} ${oirts} -g ${rg} --image ${img} --name ${name} -e METAPARTICLE_IN_CONTAINER=true`);
        setInterval( () => {
            shell.exec(`az container logs -g ${rg} -n ${name}`);
        }, 5000);
    };

    module.exports.cancel = (name) => {
        shell.exec(`az container delete --yes -g ${rg} -n ${name}`);
    };
})();
