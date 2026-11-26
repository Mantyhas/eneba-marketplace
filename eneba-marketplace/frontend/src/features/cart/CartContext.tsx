import {createContext,useContext,useEffect,useState,type ReactNode} from 'react';
export interface CartItem{offerId:number;quantity:number}
export interface Line extends CartItem{name:string;slug:string;image:string;seller:string;priceCents:number;stock:number}
export interface Quote{lines:Line[];totalCents:number}
interface Cart{items:CartItem[];add:(id:number)=>void;change:(id:number,qty:number)=>void;clear:()=>void}
const Context=createContext<Cart|null>(null);
function initial():CartItem[]{try{const v=JSON.parse(localStorage.getItem('cart')||'[]');return Array.isArray(v)?v.filter(i=>Number.isInteger(i.offerId)&&i.offerId>0&&Number.isInteger(i.quantity)&&i.quantity>0&&i.quantity<=10):[]}catch{return []}}
export function CartProvider({children}:{children:ReactNode}){const [items,setItems]=useState(initial);useEffect(()=>{localStorage.setItem('cart',JSON.stringify(items))},[items]);function update(fn:(old:CartItem[])=>CartItem[]){setItems(old=>fn(old))}function add(id:number){update(old=>old.some(i=>i.offerId===id)?old.map(i=>i.offerId===id?{...i,quantity:Math.min(10,i.quantity+1)}:i):[...old,{offerId:id,quantity:1}])}function change(id:number,qty:number){update(old=>qty<=0?old.filter(i=>i.offerId!==id):old.map(i=>i.offerId===id?{...i,quantity:Math.min(10,qty)}:i))}return <Context.Provider value={{items,add,change,clear:()=>update(()=>[])}}>{children}</Context.Provider>}
export function useCart(){const c=useContext(Context);if(!c)throw new Error('CartProvider missing');return c}
