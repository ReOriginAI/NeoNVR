const {test}=require('node:test');
const assert=require('node:assert/strict');
const fs=require('node:fs');
const vm=require('node:vm');

test('saved ranges render immediately and refresh to new analysis',async()=>{
 const elements=new Map();
 for(const id of ['mainVideo','mainSeek','detectionStatus','cameraMotionLegend','cameraMotionStatus','cameraIntensityLegend']) {
  elements.set(id,{classList:{toggle(){},remove(){}},addEventListener(){},textContent:'',title:'',currentTime:1.5});
 }
 const ranges=[];
 elements.get('mainSeek').setDetections=result=>ranges.push(result);
 elements.get('mainSeek').setCameraMotion=()=>{};
 const timers=[];
 let saved=true;
 const result={duration:4,intervals:[{kind:'person',start:1,end:2}],regions:[{start:0,end:4}],motion_frames:8,frames:8};
 const context=vm.createContext({document:{getElementById:id=>elements.get(id)},
  window:{addEventListener(){}},AbortController,URLSearchParams,
  setTimeout:fn=>{timers.push(fn);return timers.length},clearTimeout(){},
  fetch:async url=>({ok:true,json:async()=>url.startsWith('/api/detections')?{status:'complete',detections:result,stale:saved,message:'Saved settings changed'}:{motion:null}})});
 vm.runInContext(fs.readFileSync('static/detection-viewer.js','utf8'),context);
 context.window.NeoNVRDetections.loadClip({camera:'Garden',file:'clip.mp4'});
 await new Promise(resolve=>setImmediate(resolve));
 assert.equal(ranges.at(-1),result);
 assert.match(elements.get('detectionStatus').textContent,/Saved analysis/);
 assert.equal(timers.length,1);
 saved=false;
 await timers[0]();
 assert.equal(ranges.at(-1),result);
 assert.doesNotMatch(elements.get('detectionStatus').textContent,/Saved analysis/);
 assert.equal(timers.length,1);
});
