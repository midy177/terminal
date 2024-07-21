export namespace api {
	
	export class CPU {
	    user: number;
	    nice: number;
	    system: number;
	    idle: number;
	    ioWait: number;
	    irq: number;
	    softIrq: number;
	    steal: number;
	    guest: number;
	
	    static createFrom(source: any = {}) {
	        return new CPU(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user = source["user"];
	        this.nice = source["nice"];
	        this.system = source["system"];
	        this.idle = source["idle"];
	        this.ioWait = source["ioWait"];
	        this.irq = source["irq"];
	        this.softIrq = source["softIrq"];
	        this.steal = source["steal"];
	        this.guest = source["guest"];
	    }
	}
	export class FileSystem {
	    mountPoint: string;
	    used: number;
	    free: number;
	
	    static createFrom(source: any = {}) {
	        return new FileSystem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mountPoint = source["mountPoint"];
	        this.used = source["used"];
	        this.free = source["free"];
	    }
	}
	export class Network {
	    ipv4: string;
	    ipv6: string;
	    rx: number;
	    tx: number;
	
	    static createFrom(source: any = {}) {
	        return new Network(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ipv4 = source["ipv4"];
	        this.ipv6 = source["ipv6"];
	        this.rx = source["rx"];
	        this.tx = source["tx"];
	    }
	}
	export class Stat {
	    uptime: number;
	    hostname: string;
	    load1: string;
	    load5: string;
	    load10: string;
	    runningProcess: string;
	    totalProcess: string;
	    memTotal: number;
	    memAvailable: number;
	    memFree: number;
	    memBuffers: number;
	    memCached: number;
	    swapTotal: number;
	    swapFree: number;
	    fileSystems: FileSystem[];
	    network: {[key: string]: Network};
	    cpu: CPU;
	
	    static createFrom(source: any = {}) {
	        return new Stat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uptime = source["uptime"];
	        this.hostname = source["hostname"];
	        this.load1 = source["load1"];
	        this.load5 = source["load5"];
	        this.load10 = source["load10"];
	        this.runningProcess = source["runningProcess"];
	        this.totalProcess = source["totalProcess"];
	        this.memTotal = source["memTotal"];
	        this.memAvailable = source["memAvailable"];
	        this.memFree = source["memFree"];
	        this.memBuffers = source["memBuffers"];
	        this.memCached = source["memCached"];
	        this.swapTotal = source["swapTotal"];
	        this.swapFree = source["swapFree"];
	        this.fileSystems = this.convertValues(source["fileSystems"], FileSystem);
	        this.network = this.convertValues(source["network"], Network, true);
	        this.cpu = this.convertValues(source["cpu"], CPU);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

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

