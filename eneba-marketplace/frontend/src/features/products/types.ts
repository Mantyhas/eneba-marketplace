export interface Product {id:number;slug:string;name:string;platform:string;region:string;description:string;image:string;originalCents:number;priceCents:number;popularity:number}
export interface Offer {id:number;priceCents:number;seller:{name:string;rating:number};stock:number}
