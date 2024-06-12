// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';
import {termx} from '../models';

export function AddHost(arg1:main.HostEntry):Promise<void>;

export function ClosePty(arg1:string):Promise<void>;

export function CreateLocalPty(arg1:termx.SystemShell):Promise<void>;

export function CreateSshPty(arg1:string,arg2:number,arg3:number,arg4:number):Promise<void>;

export function DelHost(arg1:number):Promise<void>;

export function GetHost(arg1:number):Promise<Array<main.HostEntry>>;

export function GetLocalPtyList():Promise<Array<termx.SystemShell>>;

export function ResizePty(arg1:string,arg2:number,arg3:number):Promise<void>;

export function UpdHost(arg1:main.HostEntry):Promise<void>;

export function WriteToPty(arg1:string,arg2:Array<number>):Promise<void>;
