const {test}=require('node:test');
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

class Element {
  constructor(){this.children=[];this.style={};this.dataset={};this.events={};this.attrs={};}
  append(...children){this.children.push(...children)}
  replaceChildren(){this.children=[]}
  setAttribute(name,value){this.attrs[name]=value}
  addEventListener(event,fn){this.events[event]=fn}
}
function fixture(){
  const registry=new Map();
  const context=vm.createContext({HTMLElement:Element,document:{createElement:()=>new Element()},
    customElements:{get:name=>registry.get(name),define:(name,cls)=>registry.set(name,cls)}});
  vm.runInContext(fs.readFileSync('static/media-chrome-lite.js','utf8'),context);
  const seek=new (registry.get('media-time-range'))();
  seek.lanes=new Element();seek.video={duration:30,currentTime:0};
  return seek;
}

test('padded MOG2 regions use curly brackets and seek with clip-relative timestamps',()=>{
  const seek=fixture();
  seek.setDetections({duration:30,regions:[{start:0,end:4},{start:10,end:16}],intervals:[{kind:'person',start:11,end:12,confidence:.9}]});
  assert.equal(seek.lanes.children.length,1);
  const track=seek.lanes.children[0];
  const roi={children:track.children.filter(marker=>marker.dataset.kind==='motion_roi')};
  assert.deepEqual(roi.children.map(marker=>marker.textContent),['{','}','{','}']);
  assert.equal(roi.children[3].style.left,(16/30*100)+'%');
  roi.children[2].events.click();
  assert.equal(seek.video.currentTime,10);
  const ai=seek.lanes.children.find(lane=>lane.dataset.kind==='combined');
  assert.deepEqual(ai.children.filter(marker=>marker.className==='ai-bracket').map(marker=>marker.textContent),['[',']']);
  seek.setDetections(null);
  assert.equal(seek.lanes.children.length,0);
  assert.equal(seek.lanes.hidden,true);
});

test('invalid ROI windows do not create markers',()=>{
  const seek=fixture();
  seek.setDetections({regions:[{start:-1,end:4},{start:4,end:3},{start:5,end:40},{start:NaN,end:2}],intervals:[]});
  assert.equal(seek.lanes.children.length,0);
});
