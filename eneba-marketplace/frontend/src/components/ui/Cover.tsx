import {useState} from 'react';
export function Cover({src,name,className=''}:{src:string;name:string;className?:string}){const [failed,setFailed]=useState(false);return failed?<div className={`cover-fallback ${className}`}>{name}</div>:<img className={className} src={src} alt={name} onError={()=>setFailed(true)} loading="lazy"/>}
