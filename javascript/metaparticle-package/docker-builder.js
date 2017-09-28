(function () {
    var shell = require('shelljs');

    module.exports.build = (img) => {
        shell.exec(`docker build -t ${img} .`);
    };

    module.exports.publish = (img) => {
        shell.exec(`docker push ${img}`);
    };
})();