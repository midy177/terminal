export namespace termx {
    export interface SystemShell {
        ID: string;
        Name: string;
        Command: string;
        Args: string[];
        Env: string[];
        Cwd: string;
        Icon: string;
    }
}
