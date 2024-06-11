export namespace termx {
	
	export class SystemShell {
	    id: string;
	    name: string;
	    command: string;
	    args: string[];
	    env: string[];
	    cwd: string;
	    icon: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemShell(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.command = source["command"];
	        this.args = source["args"];
	        this.env = source["env"];
	        this.cwd = source["cwd"];
	        this.icon = source["icon"];
	    }
	}

}

