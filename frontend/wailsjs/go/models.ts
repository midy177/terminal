export namespace logic {
	
	export class FileInfo {
	    name: string;
	    full_path: string;
	    size: string;
	    mode: string;
	    mod_time: number;
	    is_dir: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.full_path = source["full_path"];
	        this.size = source["size"];
	        this.mode = source["mode"];
	        this.mod_time = source["mod_time"];
	        this.is_dir = source["is_dir"];
	    }
	}
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
	    id: number;
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

