# Species pose profile draft

Purpose: make the batch-prompt workflow ask the right motion questions for each animal before generating a sheet.

This is a prompt-design artifact, not accepted source art.

## Generic Questions

Before any multi-frame prompt, define:

- Which body parts may move?
- Which body parts must stay locked for source-family identity?
- Which motion primitive is being tested: idle, walk, scurry, forage, turn, eat, ground check, rest, groom, or reaction?
- Are feet/limbs allowed to separate from the body silhouette, or must they remain attached filled shapes?
- Which angles or phase values should change per frame?
- Which anatomy is most likely to break for this species?
- Which words are banned because they cause artifacts for this species?

## Chinchilla Standard Gray

Move:

- small attached front-foot oval angle
- small attached rear-foot oval angle
- tiny body vertical offset
- curled tail counter-angle
- tiny ear tilt for idle only

Lock:

- compact round body ratio
- large rounded ears
- curled bushy tail attached to body
- cream belly shape
- clean black eye and small closed muzzle

Risk:

- foot/belly pinholes
- tail-base gaps
- detached whisker or foot specks
- shelf-like underside
- body drifting into upright mouse/degu shape

Prompt guard:

- feet are small filled oval pads attached to the underside
- tail is one continuous attached curled shape
- no lifted feet, no open gap under belly, no separated toes

## Momonga / Sugar Glider

Move:

- membrane spread/compression phase
- tail sweep angle
- forelimb reach angle
- hindlimb tucked/extended phase

Lock:

- readable gliding membrane
- broad tail, not thin mouse tail
- side-view compact body
- large dark eye

Risk:

- membrane disappearing
- tail becoming thin
- animal turning into a mouse/squirrel
- limb membranes becoming detached chips

Prompt guard:

- membrane is one continuous filled shape connecting forelimb and hindlimb
- no detached membrane fragments
- keep broad tail attached

## Golden Syrian Hamster

Move:

- tiny forepaw tuck/reach
- cheek/body bob
- short foot pad angle
- whisker/muzzle nibble phase for forage/eat

Lock:

- round cheeks
- compact body
- no visible long tail
- tiny rounded ears

Risk:

- long degu-like tail
- rat/mouse elongation
- detached foot specks
- mouth/cheek collapse

Prompt guard:

- no tail or only invisible nub
- cheeks remain round and filled
- feet stay small attached filled pads

## Rabbit Chestnut Agouti

Move:

- hind-foot compression/extension angle
- ear tilt angle
- body height offset
- forepaw tuck

Lock:

- long ears fully inside frame
- rabbit hind-foot silhouette
- compact side-view body
- chestnut agouti coat family

Risk:

- cropped ears
- detached hind foot
- floor/shadow under hop
- body stretching into dog/cat shape

Prompt guard:

- ears fully inside cell with green padding
- hind feet are attached filled shapes
- no floor, no landing shadow, no dust

## Gecko Gray-Brown

Move:

- front toe pad angle
- rear toe pad angle
- tail curve phase
- low body crawl offset
- head tilt for reaction

Lock:

- low crawler body height remains readable
- thick attached tail
- four visible feet/toe pads when needed
- gecko head shape

Risk:

- snake-like thin strip
- missing toes
- detached toe chips
- tail disconnection
- black/red face or mouth collapse

Prompt guard:

- toes are short attached pads, not separated fragments
- tail is thick and attached
- body remains low but not too thin

