// Run with node --test tests/recordings-ui.test.cjs (no dependencies).
const {test} = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');

function fixture() {
  const elements = new Map();
  function element() {
    const events = new Map();
    return {value:'',disabled:false, dataset:{}, style:{}, attrs:{}, children:[], classList:{toggle(){},add(){},contains(){return false}},
      setAttribute(name,value){this.attrs[name]=value},removeAttribute(name){delete this.attrs[name]},
      append(...children){this.children.push(...children)},appendChild(child){this.children.push(child)},
      getBoundingClientRect(){return {left:20,top:20,bottom:64}},
      load(){this.loads=(this.loads||0)+1},play(){this.paused=false;return Promise.resolve()},pause(){this.paused=true},
      offsetWidth:220,offsetHeight:144,showModal(){this.open=true},close(){this.open=false},remove(){this.removed=true},
      focus(){},
      addEventListener(type, fn, capture) {
        const handlers = events.get(type) || [];
        if (capture === true) handlers.unshift(fn); else handlers.push(fn);
        events.set(type, handlers);
      },
      emit(type, extra={}) {
        let stopped = false;
        const evt = {pointerType:'touch',isPrimary:true,clientX:0,clientY:0,
          target:{closest:()=>null},composedPath:()=>[],defaultPrevented:false,preventDefault(){this.defaultPrevented=true},stopPropagation(){},stopImmediatePropagation(){stopped=true},...extra};
        for (const fn of events.get(type) || []) {fn(evt); if(stopped)break;}
      }
    };
  }
  let timer;
  const document = {...element(), createElement:()=>element(), getElementById(id) {
    if(!elements.has(id)) elements.set(id,element());
    return elements.get(id);
  }, querySelectorAll:()=>[],body:{...element(),classList:{toggle(){}}}};
  const context = vm.createContext({document,window:{innerWidth:1024,innerHeight:768},URLSearchParams,setTimeout:fn=>{timer=fn; return 1},clearTimeout:()=>{timer=null}});
  const source = fs.readFileSync('main.go','utf8').split('<script>')[1].split('</script>')[0]
    .replace('installVideoPanZoom();','').replace('installSideRailScrollCapture();','')
    .replace('setZoomMode(false);','').replace('init();','');
  vm.runInContext(source,context);
  vm.runInContext("visibleClips = [0,1,2,3].map(i => ({camera:'Garden',file:String(i)})); activeKey = keyFor(visibleClips[0]);",context);
  return {run:code=>vm.runInContext(code,context),context,elements,element,hold(){assert.ok(timer);const fn=timer;timer=null;fn()}};
}

test('navigation reflects visible list boundaries and follows list order',()=>{
  const f=fixture();
  f.run('updateVideoNavigation()');
  assert.equal(f.elements.get('previousVideo').disabled,true);
  assert.equal(f.elements.get('nextVideo').disabled,false);
  f.run('loadMainPlayer = clip => {activeKey = keyFor(clip); updateVideoNavigation()}; playRelative(1)');
  assert.equal(f.run('activeKey'),'Garden/1');
  f.run('playRelative(1); playRelative(1)');
  assert.equal(f.elements.get('nextVideo').disabled,true);
  f.run('playRelative(-1)');
  assert.equal(f.run('activeKey'),'Garden/2');
  f.run('visibleClips = []; updateVideoNavigation()');
  assert.equal(f.elements.get('previousVideo').disabled,true);
  assert.equal(f.elements.get('nextVideo').disabled,true);
});

test('touch hold enters selection, then extends range and suppresses click',()=>{
  const f=fixture();const first=f.element(),last=f.element();
  f.context.first=first;f.context.last=last;
  f.run('installHoldSelection(first,0); installHoldSelection(last,3)');
  let clicks=0;last.addEventListener('click',()=>clicks++);
  first.emit('pointerdown');f.hold();first.emit('pointerup');first.emit('click');
  assert.equal(f.run('selectMode'),true);
  assert.equal(f.run('selectedKeys.size'),1);
  last.emit('pointerdown');f.hold();last.emit('pointerup');last.emit('click');
  assert.equal(f.run('selectedKeys.size'),4);
  assert.equal(clicks,0);
  last.emit('click');assert.equal(clicks,1);
});

test('scrolling and cancelled touch do not select clips; mouse still supports shift',()=>{
  const f=fixture();const el=f.element();f.context.el=el;f.run('installHoldSelection(el,0)');
  el.emit('pointerdown');el.emit('pointermove',{clientY:20});assert.throws(()=>f.hold());
  el.emit('pointerdown');el.emit('pointercancel');assert.throws(()=>f.hold());
  el.emit('pointerdown',{pointerType:'mouse'});assert.throws(()=>f.hold());
  assert.equal(f.run('selectedKeys.size'),0);
  f.run('toggleSelection(0,false); toggleSelection(2,true)');
  assert.equal(f.run('selectedKeys.size'),3);
});


test('Quarantine tab fetches all dates and restores Clips date filtering',async()=>{
  const f=fixture();
  f.run(`
    var requestedURL = '';
    renderMarks = () => {};
    renderClips = () => {};
    renderDatePills = () => {};
    getJSON = async url => {requestedURL = url; return {clips:[{camera:'Garden',file:'.neonvr_pruned/weak.mp4'}]}};
    currentView = 'quarantine';
  `);
  await f.run('refreshClips(false)');
  assert.equal(f.run('requestedURL'),'/api/quarantine?');
  assert.equal(f.run('allClips[0].file'),'.neonvr_pruned/weak.mp4');
  f.run('refreshClips = async () => {}; refreshDates = async () => {}; setView("quarantine")');
  assert.equal(f.run('currentView'),'quarantine');
  assert.equal(f.elements.get('date').disabled,true);
  f.run('setView("clips")');
  assert.equal(f.elements.get('date').disabled,false);
});


test('All tab fetches both locations for the selected date',async()=>{
  const f=fixture();
  f.run(`renderMarks=()=>{}; renderClips=()=>{}; renderDatePills=()=>{};
    var requestedURL='';
    getJSON=async url=>{requestedURL=url; return {clips:[]}};
    currentView='all'; dateEl.value='2026-10-01';`);
  await f.run('refreshClips(false)');
  assert.equal(f.run('requestedURL'),'/api/all?date=2026-10-01');
});

test('active video scrolls only the rail and never calls scrollIntoView',()=>{
  const f=fixture();
  const rail=f.elements.get('sideRail');rail.scrollTop=0;
  rail.getBoundingClientRect=()=>({top:100,bottom:500});
  const active={dataset:{key:'Garden/0'},classList:{toggle(){},contains:name=>name==='sideThumb'},
    getBoundingClientRect:()=>({top:530-rail.scrollTop,bottom:620-rail.scrollTop}),
    scrollIntoView(){assert.fail('scrollIntoView scrolls page ancestors')}};
  f.context.document.querySelectorAll=()=>[active];
  f.run('syncActiveCard(true)');
  assert.equal(rail.scrollTop,120);
  f.run('syncActiveCard(true)');
  assert.equal(rail.scrollTop,120);
});

test('Shift arrows navigate while a button is focused and ignore text inputs',()=>{
  const f=fixture();
  f.run('loadMainPlayer=clip=>{activeKey=keyFor(clip);updateVideoNavigation()}');
  f.context.document.emit('keydown',{key:'ArrowRight',shiftKey:true,composedPath:()=>[{tagName:'BUTTON'}]});
  assert.equal(f.run('activeKey'),'Garden/1');
  f.context.document.emit('keydown',{key:'ArrowLeft',shiftKey:true,composedPath:()=>[{tagName:'BUTTON'}]});
  assert.equal(f.run('activeKey'),'Garden/0');
  f.context.document.emit('keydown',{key:'ArrowRight',shiftKey:true,composedPath:()=>[{tagName:'INPUT'}]});
  assert.equal(f.run('activeKey'),'Garden/0');
});

test('context menu offers download, delete and the correct location action',()=>{
  const f=fixture();const clip={camera:'Garden',file:'example.mp4',quarantined:false};
  f.context.clip=clip; f.context.target=f.element();
  f.run('installClipMenu(target,clip)');
  f.context.target.emit('contextmenu',{pointerType:'mouse',clientX:1000,clientY:750});
  let menu=f.run('clipMenu');
  assert.deepEqual(menu.children.map(el=>el.textContent),['Download clip','Quarantine','Delete']);
  assert.equal(menu.children[0].href,'/download?camera=Garden&file=example.mp4');
  assert.equal(menu.style.left,'796px');assert.equal(menu.style.top,'616px');
  f.run('clip.quarantined=true; openClipMenu(clip,20,20,target)');
  assert.equal(menu.removed,true);
  assert.equal(f.run('clipMenu').children[1].textContent,'Recover from quarantine');
  f.run('closeClipMenu()');
  f.context.target.emit('contextmenu',{pointerType:'touch'});
  assert.equal(f.run('clipMenu'),null);
});


test('sparkline point filters ten minutes either side and double-click restores all day',()=>{
  const f=fixture();
  f.run(`renderClips=()=>{};
    timeCenterEl.value=43200;
    loadMainPlayer=()=>{};
    allClips=[{start_seconds:42599},{start_seconds:42600},{start_seconds:43200},{start_seconds:43800},{start_seconds:43801}];
    updateRangeVisuals();`);
  assert.equal(f.run('allClips.filter(clipMatchesTime).length'),5);
  f.elements.get('timeline').emit('dblclick');
  assert.equal(f.run('timeFilterEnabled'),true);
  assert.equal(f.run('allClips.filter(clipMatchesTime).length'),3);
  assert.equal(f.elements.get('rangeFill').hidden,false);
  f.elements.get('timeline').emit('dblclick');
  assert.equal(f.run('timeFilterEnabled'),false);
  assert.equal(f.run('allClips.filter(clipMatchesTime).length'),5);
  f.elements.get('timeFilterToggle').emit('click');
  assert.equal(f.run('timeFilterEnabled'),true);
  f.elements.get('timeCenter').value=43800;
  f.elements.get('timeCenter').emit('input');
  assert.equal(f.run('allClips.filter(clipMatchesTime).length'),3);
});

test('sparkline window clips at midnight and activity counts quarantine separately',()=>{
  const f=fixture();
  f.run('timeCenterEl.value=0');
  assert.equal(f.run('selectedTimeWindow().start'),0);
  assert.equal(f.run('selectedTimeWindow().end'),600);
  f.run('timeCenterEl.value=86399');
  assert.equal(f.run('selectedTimeWindow().start'),85740);
  assert.equal(f.run('selectedTimeWindow().end'),86399);
  f.run(`allClips=[{start_seconds:0},{start_seconds:59,quarantined:true},{start_seconds:86399}]; renderMarks()`);
  assert.equal(f.run('activityBins(allClips).total[0]'),2);
  assert.equal(f.run('activityBins(allClips).quarantine[0]'),1);
  assert.equal(f.run('activityBins(allClips).total[1439]'),1);
  assert.match(f.elements.get('sparklineLine').attrs.d,/^M/);
});

test('All uses date pills and includes quarantine dates when loading dates',async()=>{
  const f=fixture();
  f.run(`currentView='all'; var requestedURL='';
    renderDatePills=()=>{};
    getJSON=async url=>{requestedURL=url;return {dates:['2026-10-02','2026-10-01']}};`);
  await f.run('refreshDates(true)');
  assert.equal(f.run('requestedURL'),'/api/dates?view=all');
  assert.equal(f.elements.get('date').value,'2026-10-02');
  f.run('refreshDates=async()=>{}; refreshClips=async()=>{}; setView("all")');
  assert.equal(f.elements.get('date').disabled,false);
});


test('selecting a video moves the sparkline to its minute',()=>{
  const f=fixture();
  f.run(`allClips=[{camera:'Garden',file:'one',start_seconds:43235,video_url:'/one',thumb_url:'/thumb'}];
    visibleClips=allClips; loadMainPlayer(allClips[0],false,false);`);
  assert.equal(f.elements.get('timeCenter').value,43200);
  assert.equal(f.run('activeKey'),'Garden/one');
});

test('sparkline chooses a nearby clip without snapping the point back to its time',()=>{
  const f=fixture();
  f.run(`allClips=[{camera:'Garden',file:'one',start_seconds:43235,video_url:'/one'},
    {camera:'Garden',file:'two',start_seconds:46855,video_url:'/two'}];
    visibleClips=allClips; timeCenterEl.value=45000;`);
  f.elements.get('timeCenter').emit('input');
  assert.equal(f.run('activeKey'),'Garden/one');
  assert.equal(f.elements.get('timeCenter').value,45000);
  assert.equal(f.elements.get('mainVideo').loads,1);
  // Input events within the same minute don't reload the video or rebuild cards.
  f.elements.get('timeCenter').value=45020;
  f.elements.get('timeCenter').emit('input');
  assert.equal(f.elements.get('mainVideo').loads,1);
});

test('minute selection stays stable among multiple clips within a minute',()=>{
  const f=fixture();
  f.run(`allClips=[{camera:'Garden',file:'a',start_seconds:43201},{camera:'Garden',file:'b',start_seconds:43259}];
    activeKey='Garden/b';`);
  assert.equal(f.run('nearestClipToMinute(allClips,43200).file'),'b');
  assert.equal(f.run('activityBins(allClips).total[720]'),2);
  assert.equal(f.run('activityBins(allClips).total.length'),1440);
});
