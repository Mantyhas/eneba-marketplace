import {createContext,useContext,useState,type ReactNode} from 'react';
import {api,post} from '../../api/client';
interface User{id:number;email:string}
interface Session{token:string;user:User}
interface Auth {session:Session|null;signIn:(email:string,password:string,register:boolean)=>Promise<void>;signOut:()=>Promise<void>;clear:()=>void}
const Context=createContext<Auth|null>(null);
function initial(){try{return JSON.parse(sessionStorage.getItem('session')||'null') as Session|null}catch{return null}}
export function AuthProvider({children}:{children:ReactNode}){const [session,setSession]=useState<Session|null>(initial);function clear(){sessionStorage.removeItem('session');setSession(null)} async function signIn(email:string,password:string,register:boolean){const s=await api<Session>(`/auth/${register?'register':'login'}`,post({email,password}));sessionStorage.setItem('session',JSON.stringify(s));setSession(s)} async function signOut(){try{if(session)await api('/auth/logout',post({}),session.token)}finally{clear()}}return <Context.Provider value={{session,signIn,signOut,clear}}>{children}</Context.Provider>}
export function useAuth(){const c=useContext(Context);if(!c)throw new Error('AuthProvider missing');return c}
