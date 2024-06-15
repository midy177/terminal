export namespace logic {
	
	
	export class HostEntry {
	    id: number;
	    is_folder: boolean;
	    label: string;
	    username: string;
	    address: string;
	    port: number;
	    password: string;
	    folder_id: number;
	    key_id: number;
	
	    static createFrom(source: any = {}) {
	        return new HostEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.is_folder = source["is_folder"];
	        this.label = source["label"];
	        this.username = source["username"];
	        this.address = source["address"];
	        this.port = source["port"];
	        this.password = source["password"];
	        this.folder_id = source["folder_id"];
	        this.key_id = source["key_id"];
	    }
	}
	export class KeyEntry {
	    id?: number;
	    label?: string;
	    content?: string;
	
	    static createFrom(source: any = {}) {
	        return new KeyEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.content = source["content"];
	    }
	}

}

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

