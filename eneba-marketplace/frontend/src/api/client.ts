const base = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';
export class ApiError extends Error { constructor(public status:number,message:string){super(message)} }
export async function api<T>(path:string, options:RequestInit={}, token?:string):Promise<T>{
 const headers = new Headers(options.headers); if(options.body) headers.set('Content-Type','application/json'); if(token) headers.set('Authorization',`Bearer ${token}`);
 const response=await fetch(`${base}${path}`,{...options,headers});
 const data=await response.json().catch(()=>null);
 if(!response.ok) throw new ApiError(response.status,data?.error || 'Request failed');
 return data as T;
}
export const post=(body:unknown):RequestInit=>({method:'POST',body:JSON.stringify(body)});
export const money=(cents:number)=>new Intl.NumberFormat('en-IE',{style:'currency',currency:'EUR'}).format(cents/100);
