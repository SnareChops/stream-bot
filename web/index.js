(() => {
  // core/engine.ts
  var Engine = class {
    constructor(context) {
      this.context = context;
    }
    prev = 0;
    game;
    runGame(game) {
      this.game = game;
      window.requestAnimationFrame(this.tick.bind(this));
    }
    tick(timestamp) {
      const delta = timestamp - this.prev;
      this.prev = timestamp;
      if (delta > 0) {
        this.game.update(delta);
      }
      this.context.reset();
      this.game.draw(this.context);
      window.requestAnimationFrame(this.tick.bind(this));
    }
  };

  // core/util.ts
  function isSet(mask, state) {
    return (mask & state) === state;
  }
  function createCanvas(w, h) {
    const canvas = document.createElement("canvas");
    canvas.width = w;
    canvas.height = h;
    return [canvas, canvas.getContext("2d")];
  }

  // core/font.ts
  var Font5;
  var Font10;
  var Font15;
  var Font20;
  async function initFonts() {
    return new Promise((resolve, reject) => {
      const png = new Image();
      png.src = `assets/Fonts.png`;
      png.addEventListener("load", (event) => {
        const image = event.target;
        const [canvas, context] = createCanvas(image.naturalWidth, image.naturalHeight);
        context.drawImage(image, 0, 0);
        Font5 = loadFontFromImage(1, 5, context);
        Font10 = loadFontFromImage(7, 10, context);
        Font15 = loadFontFromImage(18, 15, context);
        Font20 = loadFontFromImage(34, 20, context);
        resolve();
      });
    });
  }
  function loadFontFromImage(start, height, context) {
    const font = {
      height,
      data: /* @__PURE__ */ new Map()
    };
    let cursor = 0;
    for (let ascii = 32; ascii < 127; ascii++) {
      cursor += 1;
      const begin = cursor;
      while (true) {
        const [r, g, b, a] = context.getImageData(cursor, start, 1, 1).data;
        if (g == 255)
          break;
        cursor += 1;
      }
      font.data.set(ascii, {
        width: cursor - begin,
        pixels: extractPixels(context, begin, start, cursor, start + height)
      });
    }
    return font;
  }
  function extractPixels(context, x1, y1, x2, y2) {
    const result = [];
    for (let y = y1; y < y2; y++) {
      for (let x = x1; x < x2; x++) {
        const [r, g, b, a] = context.getImageData(x, y, 1, 1).data;
        result.push(r == 255);
      }
    }
    return result;
  }
  function drawStringBlock(dest, bounds, kerning, leading, lines) {
    const [x, y] = bounds.vector.vec2();
    for (let i = 0; i < lines.length; i++) {
      drawString(dest, x, y + i * (leading + lines[i][0].image.height), kerning, lines[i]);
    }
  }
  function drawString(dest, x, y, kerning, letters) {
    let cursor = 0;
    for (const letter of letters) {
      dest.drawImage(letter.image, x + cursor, y);
      cursor += letter.image.width + kerning;
    }
  }
  async function imagesForString(text, font, color) {
    const result = [];
    let line = [];
    for (let i = 0; i < text.length; i++) {
      const char = text.charCodeAt(i);
      if (char == 10) {
        result.push(line);
        line = [];
        continue;
      }
      if (char < 32 || char > 127) {
        console.log("Skipped generating image for unsupported rune", char);
        continue;
      }
      const c = font.data.get(char);
      if (!!c) {
        const pixels = pixelsForLetter(font.height, c.pixels, color);
        const bitmap = await createImageBitmap(new ImageData(pixels, c.width));
        line.push({ image: bitmap, char });
      }
    }
    if (line.length > 0) {
      result.push(line);
    }
    return result;
  }
  function pixelsForLetter(width, pixels, color) {
    const letter = new Uint8ClampedArray(pixels.length * 4);
    for (let i = 0; i < pixels.length; i++) {
      if (pixels[i]) {
        letter[i * 4] = color[0];
        letter[i * 4 + 1] = color[1];
        letter[i * 4 + 2] = color[2];
        letter[i * 4 + 3] = color[3];
      }
    }
    return letter;
  }
  function wrap(bounds, kerning, leading, lines) {
    const result = [];
    let currLine = [];
    let currWidth = -kerning;
    for (const line of lines) {
      for (const letter of line) {
        currWidth += kerning + letter.image.width;
        currLine.push(letter);
        if (currWidth > bounds.dx()) {
          for (let i = currLine.length - 1; i > 0; i--) {
            if (currLine[i].char == 32) {
              result.push(currLine.slice(0, i));
              currLine = currLine.slice(i + 1);
              currWidth = -kerning;
              for (const l of currLine) {
                currWidth += kerning + l.image.width;
              }
              break;
            }
          }
        }
      }
      if (currLine.length > 0) {
        result.push(currLine);
        currLine = [];
        currWidth = -kerning;
      }
    }
    return result;
  }

  // src/game.ts
  var Game = class {
    constructor(scene) {
      this.scene = scene;
    }
    update(delta) {
      this.scene.update(delta);
    }
    draw(context) {
      this.scene.draw(context);
    }
    loadScene(scene) {
      this.scene = scene;
    }
  };

  // core/vector.ts
  var RawVector = class {
    x = 0;
    y = 0;
    z = 0;
    vec2() {
      return [this.x, this.y];
    }
    setVec2(x, y) {
      this.x = x;
      this.y = y;
    }
    vec3() {
      return [this.x, this.y, this.z];
    }
    setVec3(x, y, z) {
      this.x = x;
      this.y = y;
      this.z = z;
    }
  };

  // core/bounds.ts
  var TOP = 0;
  var CENTER = 1;
  var BOTTOM = 2;
  var LEFT = 3;
  var RIGHT = 4;
  var BaseBounds = class {
    constructor(width, height) {
      this.width = width;
      this.height = height;
    }
    vector;
    offsetX = 0;
    offsetY = 0;
    anchorX = 0;
    anchorY = 0;
    rawPos() {
      const [ox, oy] = this.offset();
      const [vx, vy] = this.vector.vec2();
      return [vx - ox, vy - oy];
    }
    setAnchor(x, y) {
      this.anchorX = x;
      this.anchorY = y;
      switch (x) {
        case LEFT:
          this.offsetX = 0;
          break;
        case CENTER:
          this.offsetX = this.width / 2;
          break;
        case RIGHT:
          this.offsetX = this.width;
          break;
      }
      switch (y) {
        case TOP:
          this.offsetY = 0;
          break;
        case CENTER:
          this.offsetY = this.height / 2;
          break;
        case BOTTOM:
          this.offsetY = this.height;
          break;
      }
    }
    vecOf(h, v) {
      let x = 0;
      let y = 0;
      const [vx, vy] = this.vector.vec2();
      switch (h) {
        case LEFT:
          x = vx - this.offsetX;
          break;
        case CENTER:
          x = vx - this.offsetX + this.width / 2;
          break;
        case RIGHT:
          x = vx - this.offsetX + this.width;
          break;
      }
      switch (v) {
        case TOP:
          y = vy - this.offsetY;
          break;
        case CENTER:
          y = vy - this.offsetY + this.height / 2;
          break;
        case BOTTOM:
          y = vy - this.offsetY + this.height;
          break;
      }
      return [x, y];
    }
    anchor() {
      return [this.anchorX, this.anchorY];
    }
    offset() {
      return [this.offsetX, this.offsetY];
    }
    size() {
      return [this.width, this.height];
    }
    setSize(w, h) {
      this.width = w;
      this.height = h;
    }
    dx() {
      return this.width;
    }
    dy() {
      return this.height;
    }
    normalVectorOf(edge) {
      switch (edge) {
        case LEFT:
          return [-1, 0];
        case TOP:
          return [0, -1];
        case RIGHT:
          return [1, 0];
        case BOTTOM:
          return [0, 1];
        default:
          throw new Error("Invalid edge");
      }
    }
    isWithin(x, y) {
      const [x1, y1] = this.rawPos();
      if (this.width == 1 && this.height == 1) {
        return x == x1 && y == y1;
      }
      const x2 = x1 + this.width;
      const y2 = y1 + this.height;
      return x > x1 && x < x2 && y > y1 && y < y2;
    }
    doesCollide(other) {
      const [w1, h1] = this.size();
      const [x1, y1] = this.rawPos();
      const [w2, h2] = other.size();
      const [x2, y2] = other.rawPos();
      return !(x2 + w2 < x1 || x2 > x1 + w1 || y2 + h2 < y1 || y2 > y1 + h1);
    }
    collisionEdges(other) {
      const [w1, h1] = this.size();
      const [x1, y1] = this.rawPos();
      const [w2, h2] = other.size();
      const [x2, y2] = other.rawPos();
      if (x1 + w1 >= x2 && x1 < x2) {
        return [LEFT, RIGHT];
      }
      if (x1 <= x2 + w2 && x1 + w1 > x2 + w2) {
        return [RIGHT, LEFT];
      }
      if (y1 + h1 >= y2 && y1 < y2) {
        return [TOP, BOTTOM];
      }
      if (y1 <= y2 + h2 && y1 + h1 > y2 + h2) {
        return [BOTTOM, TOP];
      }
      return [0, 0];
    }
  };
  var RawBounds = class extends BaseBounds {
    constructor(width = 0, height = 0) {
      super(width, height);
      this.width = width;
      this.height = height;
      this.vector = new RawVector();
    }
  };

  // core/physics.ts
  var PhysicsSimulator = class {
    entities = [];
    colliders = [];
    addEntity(entity) {
      if (!this.entities.includes(entity)) {
        this.entities.push(entity);
      }
    }
    removeEntity(entity) {
      if (this.entities.includes(entity)) {
        this.entities = this.entities.splice(this.entities.indexOf(entity), 1);
      }
    }
    addCollider(collider) {
      if (!this.colliders.includes(collider)) {
        this.colliders.push(collider);
      }
    }
    removeCollider(collider) {
      if (this.colliders.includes(collider)) {
        this.colliders = this.colliders.splice(this.colliders.indexOf(collider), 1);
      }
    }
    update(delta) {
      for (const entity of this.entities) {
        entity.update(delta);
        for (const collider of this.colliders) {
          if (isSet(entity.collision(), collider.type()) && entity.doesCollide(collider)) {
            const [e1, e2] = entity.collisionEdges(collider);
            const [vx, vy] = entity.velocity();
            const [nx, ny] = collider.normalVectorOf(e2);
            const [nvx, nvy] = reflectVelocity(vx, vy, nx, ny);
            entity.setRawVelocity(nvx, nvy);
            break;
          }
        }
      }
    }
  };
  function reflectVelocity(vx, vy, nx, ny) {
    const dotProduct = vx * nx + vy * ny;
    return [vx - 2 * dotProduct * nx, vy - 2 * dotProduct * ny];
  }

  // core/random.ts
  function intn(n) {
    return Math.floor(Math.random() * n);
  }
  function floatn(n) {
    return Math.random() * n;
  }

  // core/trig.ts
  function pointAtAngleWithDistance(x, y, angle, dist) {
    return [
      x + dist * Math.cos(angle),
      y + dist * Math.sin(angle)
    ];
  }

  // core/particle.ts
  var ParticleEmitter = class {
    bounds = new RawBounds();
    particles = [];
    images = [];
    minVelocity = 0;
    maxVelocity = 0;
    minAngle = 0;
    maxAngle = 0;
    minLife = 0;
    maxLife = 0;
    density = 0;
    duration = 0;
    init(poolSize, images) {
      this.images = images;
      this.particles = Array(poolSize);
      for (let i = 0; i < this.particles.length; i++) {
        const idx = Math.floor(Math.random() * this.images.length);
        const image = this.images[idx];
        this.particles[i] = new Particle();
        this.particles[i].bounds.setAnchor(CENTER, CENTER);
        this.particles[i].image = image;
      }
      return this;
    }
    setVelocity(min, max) {
      this.minVelocity = min;
      this.maxVelocity = max;
    }
    setAngle(min, max) {
      this.minAngle = min;
      this.maxAngle = max;
    }
    setLife(min, max) {
      this.minLife = min;
      this.maxLife = max;
    }
    setDensity(density) {
      this.density = density;
    }
    start(duration) {
      this.duration = duration;
    }
    update(delta) {
      this.duration -= delta;
      if (this.duration <= 0) {
        this.duration = 0;
      }
      let desired = delta / 10 * this.density;
      for (let particle of this.particles) {
        if (particle.active) {
          particle.update(delta);
        } else {
          if (desired > 0 && this.duration > 0) {
            const [x, y] = this.bounds.vector.vec2();
            particle.start(
              x,
              y,
              intn(this.maxLife - this.minLife) + this.minLife,
              floatn(this.maxVelocity - this.minVelocity) + this.minVelocity,
              floatn(this.maxAngle - this.minAngle) + this.minAngle
            );
            desired -= 1;
          }
        }
      }
    }
  };
  var Particle = class {
    bounds = new RawBounds();
    image;
    active = false;
    life = 0;
    velocity = 0;
    angle = 0;
    start(x, y, life, velocity, angle) {
      this.active = true;
      this.life = life;
      this.velocity = velocity;
      this.angle = angle;
      this.bounds.vector.setVec2(x, y);
    }
    update(delta) {
      if (!this.active) {
        return;
      }
      this.life -= delta;
      if (this.life <= 0) {
        this.active = false;
        return;
      }
      let [x, y] = this.bounds.vector.vec2();
      [x, y] = pointAtAngleWithDistance(x, y, this.angle, this.velocity * delta);
      this.bounds.vector.setVec2(x, y);
    }
    Image() {
      if (this.active) {
        return this.image;
      }
      return void 0;
    }
  };

  // core/renderer.ts
  var Renderer = class {
    background = [];
    effects = [];
    screen = [];
    addToBackground(sprite) {
      if (this.background.includes(sprite))
        return;
      this.screen.push(sprite);
    }
    removeFromBackground(sprite) {
      const i = this.background.indexOf(sprite);
      if (i > -1)
        this.background.splice(i, 1);
    }
    addEffect(effect) {
      if (this.effects.includes(effect))
        return;
      this.effects.push(effect);
    }
    removeEffect(effect) {
      const i = this.effects.indexOf(effect);
      if (i > -1)
        this.effects.splice(i, 1);
    }
    addToScreen(sprite) {
      if (this.screen.includes(sprite))
        return;
      this.screen.push(sprite);
    }
    removeFromScreen(sprite) {
      const i = this.screen.indexOf(sprite);
      if (i > -1)
        this.screen.splice(i, 1);
    }
    draw(ctx2) {
      for (let effect of this.effects) {
        for (let particle of effect.particles) {
          const image = particle.Image();
          if (!!image) {
            ctx2.drawImage(image, ...particle.bounds.rawPos());
          }
        }
      }
      this.screen.sort((a, b) => {
        const [ax, ay, az] = a.vector.vec3();
        const [bx, by, bz] = b.vector.vec3();
        return az - bz;
      });
      for (const item of this.screen) {
        const image = item.Image();
        if (!!image) {
          ctx2.drawImage(image, ...item.rawPos());
        }
      }
    }
  };

  // src/cursor.ts
  var Cursor = class {
    clicked = false;
    x = 0;
    y = 0;
    constructor() {
      document.addEventListener("click", (event) => {
        this.clicked = true;
        this.x = event.clientX;
        this.y = event.clientY;
      });
    }
    update() {
      this.clicked = false;
    }
  };

  // src/socket.ts
  var Socket = class {
    #url;
    #socket;
    #handlers = [];
    #interval = 1e3;
    #maxInterval = 3e4;
    #decay = 2;
    constructor(url) {
      this.#url = url;
      this.#connect();
    }
    addHandler(handler) {
      this.#handlers.push(handler);
    }
    send(signal) {
      this.#socket?.send(JSON.stringify(signal));
    }
    #connect() {
      this.#socket = new WebSocket(this.#url);
      this.#socket.onopen = () => {
        console.log("Socket connection open");
        this.#interval = 1e3;
      };
      this.#socket.onclose = (event) => {
        console.log("Socket connection closed", event);
        this.#reconnect();
      };
      this.#socket.onerror = (err) => console.error("Socket connection error", err);
      this.#socket.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          for (const handler of this.#handlers) {
            handler(message);
          }
        } catch (err) {
          console.log("Failed to parse message data", err);
        }
      };
    }
    #reconnect() {
      setTimeout(() => {
        console.log(`Socket reconnecting... (delay: ${this.#interval}ms)`);
        this.#connect();
        this.#interval = Math.min(
          this.#interval * this.#decay,
          this.#maxInterval
        );
      }, this.#interval);
    }
  };

  // src/collision.ts
  var WallCollider = class extends RawBounds {
    constructor(x, y, w, h) {
      super(w, h);
      this.vector.setVec2(x, y);
    }
    type() {
      return 1 /* WALL */;
    }
  };

  // src/sprites/alert.ts
  var AlertBox = class extends RawBounds {
    constructor(text, rate, wait) {
      super(0, 0);
      this.text = text;
      this.rate = rate;
      this.wait = wait;
      this.textBounds.vector.setVec2(22, 15);
      this.box = new Image();
      this.box.src = "assets/AlertBox.png";
      this.box.addEventListener("load", ((x) => {
        const target = x.target;
        this.setSize(target.naturalWidth, target.naturalHeight);
        [this.canvas, this.context] = createCanvas(...this.size());
        this.render();
      }).bind(this));
      imagesForString(this.text, Font15, [1, 249, 255, 255]).then((lines) => {
        this.lines = wrap(this.textBounds, this.kerning, this.leading, lines);
        this.render();
      });
    }
    box;
    canvas;
    context;
    cooldown = 0;
    status = 0 /* PENDING */;
    index = 0;
    lines;
    kerning = 2;
    leading = 4;
    textBounds = new RawBounds(377, 119);
    done;
    start(done) {
      this.cooldown = this.rate;
      this.index = 0;
      this.status = 1 /* ANIMATING */;
      this.done = done;
    }
    update(delta) {
      if (this.status === 1 /* ANIMATING */) {
        this.cooldown -= delta;
        if (this.cooldown < 1) {
          this.cooldown = this.rate + this.cooldown;
          this.index += 1;
          if (this.index === this.text.length) {
            this.status = 2 /* WAITING */;
          }
          this.render();
          return;
        }
      }
      if (this.status === 2 /* WAITING */) {
        this.wait -= delta;
        if (this.wait < 1) {
          this.status = 3 /* DONE */;
          this.done();
        }
      }
    }
    Image() {
      return this.canvas;
    }
    render() {
      if (!this.context || !this.lines)
        return;
      this.context.reset();
      this.context.drawImage(this.box, 0, 0);
      if (this.index === 0)
        return;
      let count = 0;
      const linesToDraw = [];
      for (const line of this.lines) {
        if (line.length + count < this.index) {
          count += line.length;
          linesToDraw.push(line);
        } else {
          const end = this.index - count;
          linesToDraw.push(line.slice(0, end));
          break;
        }
      }
      drawStringBlock(this.context, this.textBounds, this.kerning, this.leading, linesToDraw);
    }
  };

  // src/alert.ts
  function FollowAlert(signal) {
    const follow = signal;
    return new AlertBox(`${follow.username} followed!`, 50, 5e3);
  }
  function CheerAlert(signal) {
    const cheer = signal;
    return new AlertBox(`${cheer.username} wasted ${cheer.amount} bits on this dumb little box${!!cheer.message ? ` to say:
${cheer.message}` : "."}`, 50, 5e3);
  }
  function SubAlert(signal) {
    const sub = signal;
    return new AlertBox(`${sub.username} made a tier ${sub.tier} mistake wasting money on emotes.`, 50, 5e3);
  }
  function ResubAlert(signal) {
    const resub = signal;
    return new AlertBox(`${resub.username} has wasted tier ${resub.tier} money for ${resub.months} months.`, 50, 5e3);
  }
  function GiftAlert(signal) {
    const gift = signal;
    return new AlertBox(`${gift.anon ? "Anonymous" : gift.username} gifted a tier ${gift.tier} subscription to `, 50, 5e3);
  }

  // src/audio.ts
  async function playAudio(signal) {
    const audio = document.createElement("audio");
    audio.setAttribute("src", signal.url);
    audio.setAttribute("autoplay", "true");
    audio.setAttribute("preload", "auto");
    audio.addEventListener("ended", function() {
      this.remove();
    }, true);
    document.body.append(audio);
    try {
      await audio.play();
    } catch (err) {
      console.error("Audio failed to play", signal, err);
    }
  }

  // src/sprites/avatar.ts
  var Avatar = class _Avatar extends RawBounds {
    constructor(width, height, images) {
      super(width, height);
      this.images = images;
    }
    state = 1 /* MOUTH_CLOSED */ | 8 /* EYES_OPEN */;
    blink = 0;
    static async load() {
      const images = await Promise.all([
        loadImage("assets/gremlin/idle.png"),
        loadImage("assets/gremlin/talk.png"),
        loadImage("assets/gremlin/xd.png"),
        loadImage("assets/gremlin/xdTalk.png"),
        loadImage("assets/gremlin/wheeze.png"),
        loadImage("assets/gremlin/pog.png")
      ]);
      const scaled = images.map((x) => scaleImage(x, 250, 335));
      console.log(scaled);
      const map = /* @__PURE__ */ new Map([
        [1 /* MOUTH_CLOSED */ | 8 /* EYES_OPEN */, scaled[0]],
        [2 /* MOUTH_OPEN */ | 8 /* EYES_OPEN */, scaled[1]],
        [1 /* MOUTH_CLOSED */ | 4 /* EYES_CLOSED */, scaled[2]],
        [2 /* MOUTH_OPEN */ | 4 /* EYES_CLOSED */, scaled[3]],
        [16 /* WHEEZE */, scaled[4]],
        [32 /* POG */, scaled[5]]
      ]);
      const { width, height } = scaled.values().next().value;
      return new _Avatar(width, height, map);
    }
    signal(signal) {
      switch (signal.type) {
        case "talking":
          if (signal.talking) {
            this.state &= ~1 /* MOUTH_CLOSED */;
            this.state |= 2 /* MOUTH_OPEN */;
          } else {
            this.state &= ~2 /* MOUTH_OPEN */;
            this.state |= 1 /* MOUTH_CLOSED */;
          }
          return;
      }
    }
    update(delta) {
      this.blink -= delta;
      if (this.blink <= 0) {
        if (isSet(this.state, 4 /* EYES_CLOSED */)) {
          this.state &= ~4 /* EYES_CLOSED */;
          this.state |= 8 /* EYES_OPEN */;
          this.blink = Math.random() * 19e3 + 100;
        } else {
          this.state &= ~8 /* EYES_OPEN */;
          this.state |= 4 /* EYES_CLOSED */;
          this.blink = Math.random() * 750 + 50;
        }
      }
    }
    Image() {
      if (isSet(this.state, 16 /* WHEEZE */))
        return this.images.get(16 /* WHEEZE */);
      if (isSet(this.state, 32 /* POG */))
        return this.images.get(32 /* POG */);
      return this.images.get(this.state);
    }
  };
  async function loadImage(url) {
    return new Promise((resolve) => {
      const image = new Image();
      image.src = url;
      image.addEventListener("load", function() {
        resolve(this);
      });
    });
  }
  function scaleImage(image, maxWidth, maxHeight) {
    const scaleWidth = maxWidth / image.width;
    const scaleHeight = maxHeight / image.height;
    const scale = Math.min(scaleWidth, scaleHeight);
    const width = image.width * scale;
    const height = image.height * scale;
    const [canvas, context] = createCanvas(width, height);
    context.drawImage(image, 0, 0, width, height);
    return canvas;
  }

  // src/video.ts
  async function playVideo(signal) {
    const video = document.createElement("video");
    video.setAttribute("width", "320");
    video.setAttribute("height", "240");
    const source = document.createElement("source");
    source.setAttribute("src", signal.url);
    source.setAttribute("type", "video/mp4");
    video.append(source);
    document.body.append(video);
    video.addEventListener("ended", function() {
      this.remove();
    });
    try {
      await video.play();
    } catch (err) {
      console.error("Video failed to play", signal, err);
    }
  }

  // src/scene.ts
  var AlertScene = class {
    socket = new Socket("ws://localhost:3000/ws");
    renderer = new Renderer();
    physics = new PhysicsSimulator();
    emitter = new ParticleEmitter();
    cursor = new Cursor();
    avatar;
    walls = [];
    queue = [];
    active = [];
    init() {
      const [canvas, ctx2] = createCanvas(4, 4);
      ctx2.fillStyle = "red";
      ctx2.fillRect(0, 0, 4, 4);
      this.emitter = new ParticleEmitter().init(1e3, [canvas]);
      this.emitter.setLife(100, 200);
      this.emitter.setVelocity(0.2, 0.5);
      this.emitter.setDensity(30);
      this.emitter.setAngle(0, 2 * Math.PI);
      this.renderer.addEffect(this.emitter);
      this.walls.push(new WallCollider(-10, 0, 10, 1080));
      this.walls.push(new WallCollider(0, -10, 1920, 10));
      this.walls.push(new WallCollider(1920, 0, 10, 1080));
      this.walls.push(new WallCollider(0, 1080, 1920, 10));
      for (const collider of this.walls) {
        this.physics.addCollider(collider);
      }
      Avatar.load().then((avatar) => {
        this.avatar = avatar;
        this.avatar.setAnchor(CENTER, CENTER);
        const [w, h] = this.avatar.size();
        this.avatar.vector.setVec2(1920 - w / 2, 1080 - h / 2);
        console.log(this.avatar.vector.vec2());
        this.renderer.addToScreen(avatar);
      });
      this.socket.addHandler(this.onSignal.bind(this));
      return this;
    }
    onSignal(signal) {
      switch (signal.kind) {
        case "twitch.follow":
          return this.queue.push(FollowAlert(signal));
        case "twitch.cheer":
          return this.queue.push(CheerAlert(signal));
        case "twitch.sub":
          return this.queue.push(SubAlert(signal));
        case "twitch.resub":
          return this.queue.push(ResubAlert(signal));
        case "twitch.gift":
          return this.queue.push(GiftAlert(signal));
        case "audio":
          return playAudio(signal);
        case "avatar":
          return this.avatar?.signal(signal);
        case "video.small":
          return playVideo(signal);
      }
    }
    update(delta) {
      while (true) {
        if (this.queue.length > 0) {
          const pos = findSpace(this.active, this.queue[0], 10, 10);
          if (!pos)
            break;
          const alert = this.queue.shift();
          alert.vector.setVec2(...pos);
          alert.start(() => this.active.splice(this.active.indexOf(alert), 1));
          this.active.push(alert);
        } else
          break;
      }
      if (this.cursor.clicked) {
        this.emitter.bounds.vector.setVec2(this.cursor.x, this.cursor.y);
        this.emitter.start(100);
      }
      this.physics.update(delta);
      this.emitter.update(delta);
      this.avatar?.update(delta);
      for (const alert of this.active) {
        alert.update(delta);
      }
      this.cursor.update();
    }
    draw(ctx2) {
      this.renderer.draw(ctx2);
    }
  };
  function findSpace(active, next, paddingX, paddingY) {
    const [lastX, lastY] = active.reduce((last, alert) => {
      const [x, y] = alert.vecOf(RIGHT, BOTTOM);
      if (x > last[0])
        return [x, y];
      if (y > last[1])
        return [last[0], y];
      return last;
    }, [0, 0]);
    const [w, h] = next.size();
    if (lastY + h + paddingY <= 1080)
      return [lastX + w, lastY + h + paddingY];
    if (lastX + w + paddingX <= 1920)
      return [lastX + w + paddingX, paddingY];
    return void 0;
  }

  // src/index.ts
  var $body = document.querySelector("body");
  var $canvas = document.createElement("canvas");
  var ctx = $canvas.getContext("2d");
  $canvas.width = 1920;
  $canvas.height = 1080;
  $body.append($canvas);
  initFonts().then(() => {
    new Engine(ctx).runGame(new Game(new AlertScene().init()));
  });
})();
