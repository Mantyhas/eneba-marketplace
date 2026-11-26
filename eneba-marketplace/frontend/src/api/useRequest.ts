import {useEffect,useState} from 'react';
export function useRequest<T>(load:()=>Promise<T>,key:string){
 const [state,setState]=useState<{data?:T;error?:string;loading:boolean}>({loading:true});
 useEffect(()=>{let active=true;setState({loading:true});load().then(data=>{if(active)setState({data,loading:false})}).catch((e:unknown)=>{if(active)setState({error:e instanceof Error?e.message:'Request failed',loading:false})});return()=>{active=false};},[key]);
 return state;
}
