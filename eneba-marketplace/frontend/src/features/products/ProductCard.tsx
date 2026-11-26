import {Link} from 'react-router-dom';
import {money} from '../../api/client';
import {Cover} from '../../components/ui/Cover';
import type {Product} from './types';
export function ProductCard({product:p}:{product:Product}){const discount=p.priceCents>0?Math.max(0,Math.round((1-p.priceCents/p.originalCents)*100)):0;return <Link className="product-card" to={`/product/${p.slug}`}><div className="card-cover"><Cover src={p.image} name={p.name}/></div><div className="platform-strip">◉ {p.platform}</div><div className="card-info"><h3>{p.name} ({p.platform}) Key {p.region}</h3><span className="region">{p.region}</span>{p.priceCents>0?<><small className="from">From <s>{money(p.originalCents)}</s> <b>-{discount}%</b></small><strong className="price">{money(p.priceCents)}</strong><small className="region">Instant delivery</small></>:<strong>Sold out</strong>}</div></Link>}
export function ProductGrid({products}:{products:Product[]}){return products.length?<div className="product-grid">{products.map(p=><ProductCard key={p.id} product={p}/>)}</div>:<p className="status">No products found. Try changing your filters.</p>}
