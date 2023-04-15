var c;    // global canvas
var dm;   // desktop manager
var img;  // 图片

function setup() {
  // put setup code here
  c = createCanvas(640, 480);
  // 设置容器元素
  c.parent('myContainer')

  c.mouseClicked(customMouseClicked)
  c.mouseReleased(customMouseReleased)

  // 创建对象管理器
  dm = new DesktopManager()
  // 设置刷新速率
  frameRate(30);
}

function draw() {
  clear();
  if (img) {
    image(img, 0, 0)
  }
  if (dm) {
    dm.draw();
  }
}

function customMouseClicked() {
  console.log('click')
  if (dm) {
    dm.setActive(mouseX, mouseY)
  }
}

function customMouseReleased() {
  console.log('release')
  if (dm && dm.activeObj) {
    dm.removeActive()
  }
}

function mouseDragged() {
  if (dm && dm.activeObj) {
    dm.moveActive(mouseX, mouseY)
  } else {
    dm.setActive(mouseX, mouseY)
    dm.moveActive(mouseX, mouseY)
  }
}

class DesktopManager {
  mode = 0;          // 模式：0 一班， 1 选择器模式
  objects = [];      // 对象池
  activeObj = null;  // 活动对象【可拖拽】

  constructor() {}

  setMode(mode) {
    this.mode = mode
  }

  addObj(obj) {
    console.log('dm add obj called');
    this.objects.push(obj);
  }

  removeObj(obj) {
    let idx = this.objects.indexOf(obj);
    this.objects.splice(idx, 1);
  }

  setActive(x, y) {
    if (this.mode === 0) {
      return
    }

    if (this.activeObj && this.activeObj.hovered(x, y)) {
      // true: 依然是当前元素
      //   this.moveActive(x, y)
      return
    }

    this.activeObj = null;
    for (let i = 0; i < this.objects.length; i++) {
      if (this.objects[i].hovered(x, y)) {
        console.log('found current obj is: x=%d, y=%d , idx=%d', x, y, i);
        this.activeObj = this.objects[i];
        break;
      }
    }
    // this.moveActive(x, y)
  }

  removeActive() {
    // console.log('removeActive');
    this.activeObj = null;
  }

  deleteActive() {
    if (this.activeObj) {
      let idx = this.objects.indexOf(this.activeObj);
      this.objects.splice(idx, 1);
    }
  }

  moveActive(x, y) {
    // console.log('moveActive, x=%d, y=%d', x, y);
    if (this.activeObj && this.activeObj.hovered(x, y)) {
      // true: 依然是当前元素
      this.activeObj.move(x, y)
    }
  }

  draw() {
    console.log(this.objects.length);
    // console.log('hasActiveObj=', (this.activeObj !== null), 'mode=',
    // this.mode)
    for (let i = 0; i < this.objects.length; i++) {
      this.objects[i].draw()
    }
  }
}


class Circle {
  x = 0;
  y = 0;
  r = 10;

  constructor(x = 100, y = 100, r = 80) {
    this.x = x;
    this.y = y;
    this.r = r;
  }

  move(x, y) {
    let offsetX = Math.ceil(x - this.x);
    let offsetY = Math.ceil(y - this.y);
    this.x += offsetX;
    this.y += offsetY;
  }

  hovered(x, y) {
    let offsetX = Math.ceil(Math.abs(x - this.x))
    let offsetY = Math.ceil(Math.abs(y - this.y))
    let hovered =
        (Math.pow(offsetX, 2) + Math.pow(offsetY, 2)) <= Math.pow(this.r, 2)
    // console.log(offsetX, offsetY, this.r, hovered)
    return hovered
  }

  draw(dm) {
    // console.log(this.x, this.y, this.r, typeof (this.x))
    let {x, y, r} = this;
    fill(51);
    ellipse(x, y, r, r);
  }
}

function addCircle() {
  console.log('addCircle clicked')
  let obj = new Circle();
  if (dm) {
    dm.addObj(obj);
  }
}

function imageChange(evt) {
  // console.log(evt, evt.target);
  let file = evt.target.files[0];
  const dir = window.URL.createObjectURL(file);
  console.log(dir);

  loadImage(dir, image => {
    img = image;
    img.resize(640, 480)
    // console.log(img)
  });

  console.log('image loaded');
}

function outputImage() {
  if (c) {
    saveCanvas(c, 'output', 'jpg')
  }
}

var mode = 0;
function changeMode() {
  mode = (mode + 1) % 2;
  dm.setMode(mode)
}

function deleteObject() {
  console.log('deleteObject called');
  if (dm) {
    dm.deleteActive()
  }
}