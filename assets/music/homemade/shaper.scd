(
// Only when using Shaper
~tf = Env([-0.8, 0, 0.8], [1, 1], [8, -8]).asSignal(1024 + 1);
~tf = ~tf.asWavetableNoWrap;
~tfBuf = Buffer.loadCollection(s, ~tf);
)

{Shaper.ar(~tfBuf, SinOsc.ar(mul: 1!2))}.play