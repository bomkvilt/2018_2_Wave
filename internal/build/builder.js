const { execSync }  = require('child_process');
const path  = require('path');
const util  = require('util');
const fs    = require('fs');


// execute the command in a sync mode
function Exec(command)
{
    console.log("\t<< " + command)
    try { execSync(command, {stdio: 'inherit'}); }
    catch(e) {}
}


class Builder
{
    constructor()
    {
        console.log("-- cwd: " + process.cwd())
        let conf = path.join("apps", "apps.json")
        let rawd = fs.readFileSync(conf)
        let json = JSON.parse(rawd)
        this.configs = json
    }

    Build() 
    {
        for (let config of this.configs)
        {
            console.log('-- app: "' + config.name + '"');
            let name = config.name;
            let smid = config.smid;
            let root = path.join("apps", name)
            for (let app of config.apps)
            {
                console.log('\t{' + app.name + '}');
                let file  = path.join(root, "docker", "." + app.name + ".dockerfile")
                let image = name + "." + app.name;
                // stop a previous build
                Exec(util.format(`docker stop %s`, image));
                // build a new image
                Exec(util.format(`docker build --rm --build-arg ROOT="%s" -f "%s" -t "%s" .`, root, file , image));
                // run the image
                Exec(util.format(`docker run --rm -d --name "%s" -t "%s"`, image, image));
            }
        }
    }
}

console.log('-- building started...');

(new Builder()).Build()

console.log('-- building finished...');

process.exit(0);
