(function () {
    var shell = require('shelljs');

    portString = (options) => {
        if (!options || !options.ports) {
            return '';
        }
        var portArr = options.ports.map((port) => {
            return `-p ${port}:${port}`
        });
        return portArr.join(' ');
    };

    module.exports.run = (img, name, options) => {
        ports = portString(options);
        
        shell.exec(`docker run --name ${name} ${ports} -d ${img}`);
        shell.exec(`docker logs -f ${name}`, () => {
            console.log('done');
        });
    };

    module.exports.cancel = (name) => {
        shell.exec(`docker kill ${name}`);
        shell.exec(`docker rm ${name}`);
    };
})();