<script setup lang="ts">
import { computed, type CSSProperties } from 'vue'

export type PetMode =
  | 'idle'
  | 'walk'
  | 'zoomies'
  | 'sit'
  | 'sleep'
  | 'play'
  | 'scruff'
  | 'thrown'
  | 'dizzy'

export type PetFace =
  | 'happy'
  | 'neutral'
  | 'sad'
  | 'sleep'
  | 'surprised'
  | 'dizzy'

const props = withDefaults(
  defineProps<{
    mode: PetMode
    face: PetFace
    speed?: number
  }>(),
  { speed: 1 },
)

const speedStyle = computed(() => ({ '--pet-speed': String(props.speed) }) as CSSProperties)
</script>

<template>
  <span
    class="pet-sprite block h-full w-full"
    :class="`stage-${mode}`"
    :style="speedStyle"
  >
    <svg
      v-if="mode === 'walk' || mode === 'zoomies' || mode === 'idle'"
      viewBox="0 0 120 90"
      class="h-full w-full pose-in"
      aria-hidden="true"
    >
      <g class="rig-tail">
        <path
          d="M26 50 Q10 44 10 30 Q10 22 16 21"
          fill="none"
          stroke="#d98324"
          stroke-width="7"
          stroke-linecap="round"
        />
      </g>

      <g class="rig-leg rig-leg-bf">
        <rect x="30" y="56" width="8" height="28" rx="4" fill="#d98324" />
      </g>
      <g class="rig-leg rig-leg-ff">
        <rect x="68" y="56" width="8" height="28" rx="4" fill="#d98324" />
      </g>

      <g class="rig-body">
        <ellipse cx="55" cy="52" rx="33" ry="17" fill="#f4a23b" />
        <ellipse cx="76" cy="54" rx="9" ry="10" fill="#ffd9a0" />
        <path
          d="M40 40 q3 6 0 12"
          stroke="#d98324"
          stroke-width="3"
          stroke-linecap="round"
          fill="none"
          opacity="0.6"
        />
        <path
          d="M52 38 q3 7 0 14"
          stroke="#d98324"
          stroke-width="3"
          stroke-linecap="round"
          fill="none"
          opacity="0.6"
        />
        <path
          d="M64 40 q3 6 0 12"
          stroke="#d98324"
          stroke-width="3"
          stroke-linecap="round"
          fill="none"
          opacity="0.6"
        />
      </g>

      <g class="rig-leg rig-leg-bn">
        <rect x="36" y="58" width="8" height="26" rx="4" fill="#f4a23b" />
      </g>
      <g class="rig-leg rig-leg-fn">
        <rect x="74" y="58" width="8" height="26" rx="4" fill="#f4a23b" />
      </g>

      <g class="rig-head">
        <g class="rig-ear"><polygon points="82,24 85,8 93,20" fill="#f4a23b" /></g>
        <g class="rig-ear"><polygon points="95,20 103,8 106,24" fill="#f4a23b" /></g>
        <polygon points="84,21 86,12 91,18" fill="#ef7d8f" />
        <polygon points="97,18 102,12 104,21" fill="#ef7d8f" />
        <circle cx="93" cy="32" r="15" fill="#f4a23b" />
        <path
          d="M88 19 q2 4 0 6 M94 18 q2 4 0 6"
          stroke="#d98324"
          stroke-width="2.5"
          stroke-linecap="round"
          fill="none"
          opacity="0.6"
        />
        <ellipse cx="93" cy="38" rx="7" ry="5" fill="#ffd9a0" />

        <g v-if="face === 'sleep'" class="eyes">
          <path
            d="M85 30 h6 M97 30 h6"
            stroke="#2b2b2e"
            stroke-width="1.8"
            stroke-linecap="round"
          />
        </g>
        <g v-else-if="face === 'surprised'" class="eyes">
          <circle cx="88" cy="30" r="3.4" fill="#fff" />
          <circle cx="99" cy="30" r="3.4" fill="#fff" />
          <circle cx="88" cy="30" r="1.5" fill="#2b2b2e" />
          <circle cx="99" cy="30" r="1.5" fill="#2b2b2e" />
        </g>
        <g v-else-if="face === 'dizzy'" class="eyes">
          <path
            d="M85 27 l6 6 M91 27 l-6 6 M96 27 l6 6 M102 27 l-6 6"
            stroke="#2b2b2e"
            stroke-width="1.8"
            stroke-linecap="round"
          />
        </g>
        <g v-else-if="face === 'happy'" class="eyes">
          <path
            d="M85 30 q3 -3 6 0 M96 30 q3 -3 6 0"
            stroke="#2b2b2e"
            stroke-width="2"
            stroke-linecap="round"
            fill="none"
          />
        </g>
        <g v-else class="eyes">
          <circle cx="88" cy="30" r="2.4" fill="#2b2b2e" />
          <circle cx="99" cy="30" r="2.4" fill="#2b2b2e" />
          <circle cx="88.8" cy="29.2" r="0.8" fill="#fff" />
          <circle cx="99.8" cy="29.2" r="0.8" fill="#fff" />
        </g>

        <path d="M91.5 35 h3 l-1.5 2.2 z" fill="#e0607a" />
        <path
          v-if="face === 'dizzy'"
          d="M88 40 q2 2 4 0 q2 -2 4 0"
          stroke="#2b2b2e"
          stroke-width="1.4"
          fill="none"
          stroke-linecap="round"
        />
        <path
          v-else-if="face === 'sad'"
          d="M90 40 q3 -2 6 0"
          stroke="#2b2b2e"
          stroke-width="1.4"
          fill="none"
          stroke-linecap="round"
        />
        <path
          v-else
          d="M90 38.5 q3 2.5 6 0"
          stroke="#2b2b2e"
          stroke-width="1.4"
          fill="none"
          stroke-linecap="round"
        />
        <path
          d="M78 34 h-8 M78 37 h-8 M108 34 h8 M108 37 h8"
          stroke="#fff"
          stroke-width="1"
          stroke-linecap="round"
          opacity="0.55"
        />

      </g>
    </svg>
    <svg
      v-else-if="mode === 'scruff' || mode === 'thrown'"
      viewBox="0 0 120 90"
      class="h-full w-full pose-in"
      aria-hidden="true"
    >
      <g v-if="mode === 'thrown'" class="wind-lines">
        <path d="M15 35 Q25 35 35 35" stroke="#a1a1aa" stroke-width="2" stroke-linecap="round" fill="none" opacity="0.6"/>
        <path d="M5 55 Q20 55 40 55" stroke="#a1a1aa" stroke-width="2" stroke-linecap="round" fill="none" opacity="0.4"/>
        <path d="M85 45 Q100 45 115 45" stroke="#a1a1aa" stroke-width="2" stroke-linecap="round" fill="none" opacity="0.5"/>
      </g>

      <g class="scruff-sway">
        <g class="scruff-tail">
          <path
            d="M52 66 Q44 76 50 86"
            fill="none"
            stroke="#d98324"
            stroke-width="6"
            stroke-linecap="round"
          />
        </g>

        <ellipse cx="60" cy="50" rx="17" ry="21" fill="#f4a23b" />
        <ellipse cx="60" cy="56" rx="10" ry="12" fill="#ffd9a0" />
        <path
          d="M50 42 q3 5 0 9 M68 42 q3 5 0 9"
          stroke="#d98324"
          stroke-width="2.5"
          stroke-linecap="round"
          fill="none"
          opacity="0.6"
        />

        <g class="scruff-leg scruff-leg-bl">
          <rect x="46" y="58" width="7" height="22" rx="3.5" fill="#d98324" />
        </g>
        <g class="scruff-leg scruff-leg-br">
          <rect x="67" y="58" width="7" height="22" rx="3.5" fill="#d98324" />
        </g>
        <g class="scruff-leg scruff-leg-fl">
          <rect x="50" y="40" width="7" height="20" rx="3.5" fill="#f4a23b" />
        </g>
        <g class="scruff-leg scruff-leg-fr">
          <rect x="63" y="40" width="7" height="20" rx="3.5" fill="#f4a23b" />
        </g>

        <g class="scruff-head">
          <polygon points="48,12 45,0 56,8" fill="#f4a23b" />
          <polygon points="64,8 75,0 72,12" fill="#f4a23b" />
          <polygon points="50,10 48,4 54,8" fill="#ef7d8f" />
          <polygon points="66,8 72,4 70,10" fill="#ef7d8f" />
          <circle cx="60" cy="20" r="14" fill="#f4a23b" />
          <ellipse cx="60" cy="26" rx="6" ry="4.5" fill="#ffd9a0" />

          <g class="eyes">
            <circle cx="54" cy="18" r="3" fill="#fff" />
            <circle cx="66" cy="18" r="3" fill="#fff" />
            <circle cx="54" cy="18.5" r="1.4" fill="#2b2b2e" />
            <circle cx="66" cy="18.5" r="1.4" fill="#2b2b2e" />
          </g>

          <path d="M58.5 23 h3 l-1.5 2 z" fill="#e0607a" />
          <ellipse cx="60" cy="28" rx="1.8" ry="2.4" fill="#2b2b2e" opacity="0.8" />
          <path
            d="M46 22 h-7 M46 25 h-7 M74 22 h7 M74 25 h7"
            stroke="#fff"
            stroke-width="1"
            stroke-linecap="round"
            opacity="0.55"
          />
        </g>
      </g>
    </svg>
    <svg
      v-else-if="mode === 'sit' || mode === 'dizzy'"
      viewBox="0 0 120 90"
      class="h-full w-full pose-in"
      aria-hidden="true"
    >
      <g v-if="face === 'dizzy'" class="dizzy-stars">
        <path
          transform="translate(57 8)"
          d="M0 -3 L1 -1 L3 0 L1 1 L0 3 L-1 1 L-3 0 L-1 -1 Z"
          fill="#fbbf24"
        />
        <path
          transform="translate(71 3)"
          d="M0 -3 L1 -1 L3 0 L1 1 L0 3 L-1 1 L-3 0 L-1 -1 Z"
          fill="#fbbf24"
        />
        <path
          transform="translate(85 8)"
          d="M0 -3 L1 -1 L3 0 L1 1 L0 3 L-1 1 L-3 0 L-1 -1 Z"
          fill="#fbbf24"
        />
      </g>

      <g class="sit-body">
        <circle cx="50" cy="62" r="18" fill="#f4a23b" />
        <circle cx="42" cy="62" r="12" fill="#d98324" opacity="0.25" />
        <ellipse cx="63" cy="52" rx="15" ry="26" fill="#f4a23b" />
        <ellipse cx="68" cy="56" rx="8" ry="14" fill="#ffd9a0" />
        <path
          d="M52 34 q3 6 0 10 M58 30 q3 6 0 10"
          stroke="#d98324"
          stroke-width="3"
          stroke-linecap="round"
          fill="none"
          opacity="0.6"
        />
      </g>

      <rect x="60" y="58" width="7" height="26" rx="3.5" fill="#f4a23b" />
      <rect x="71" y="58" width="7" height="26" rx="3.5" fill="#f4a23b" />
      <ellipse cx="63.5" cy="83" rx="5" ry="3" fill="#ffd9a0" />
      <ellipse cx="74.5" cy="83" rx="5" ry="3" fill="#ffd9a0" />

      <g class="sit-tail">
        <path
          d="M33 70 Q26 85 40 85 L74 85 Q81 85 81 79"
          fill="none"
          stroke="#d98324"
          stroke-width="7"
          stroke-linecap="round"
        />
      </g>

      <g class="sit-head">
        <g class="rig-ear"><polygon points="60,14 62,0 70,10" fill="#f4a23b" /></g>
        <g class="rig-ear"><polygon points="73,10 81,0 83,14" fill="#f4a23b" /></g>
        <polygon points="62,12 64,5 68,10" fill="#ef7d8f" />
        <polygon points="75,10 79,5 81,12" fill="#ef7d8f" />
        <circle cx="71" cy="22" r="14" fill="#f4a23b" />
        <ellipse cx="71" cy="28" rx="6" ry="4.5" fill="#ffd9a0" />

        <g v-if="face === 'dizzy'" class="eyes">
          <path
            d="M63 18 l6 6 M69 18 l-6 6 M74 18 l6 6 M80 18 l-6 6"
            stroke="#2b2b2e"
            stroke-width="1.8"
            stroke-linecap="round"
          />
        </g>
        <g v-else-if="face === 'happy'" class="eyes">
          <path
            d="M63 21 q3 -3 6 0 M74 21 q3 -3 6 0"
            stroke="#2b2b2e"
            stroke-width="2"
            stroke-linecap="round"
            fill="none"
          />
        </g>
        <g v-else class="eyes">
          <circle cx="66" cy="21" r="2.2" fill="#2b2b2e" />
          <circle cx="77" cy="21" r="2.2" fill="#2b2b2e" />
          <circle cx="66.7" cy="20.3" r="0.7" fill="#fff" />
          <circle cx="77.7" cy="20.3" r="0.7" fill="#fff" />
        </g>

        <path d="M69.5 25 h3 l-1.5 2 z" fill="#e0607a" />
        <path
          v-if="face === 'dizzy'"
          d="M67 29 q2 2 4 0 q2 -2 4 0"
          stroke="#2b2b2e"
          stroke-width="1.3"
          fill="none"
          stroke-linecap="round"
        />
        <path
          v-else
          d="M68 28.5 q3 2.5 6 0"
          stroke="#2b2b2e"
          stroke-width="1.3"
          fill="none"
          stroke-linecap="round"
        />
        <path
          d="M58 25 h-7 M58 28 h-7 M84 25 h7 M84 28 h7"
          stroke="#fff"
          stroke-width="1"
          stroke-linecap="round"
          opacity="0.55"
        />
      </g>
    </svg>
    <svg
      v-else-if="mode === 'sleep'"
      viewBox="0 0 120 90"
      class="h-full w-full pose-in"
      aria-hidden="true"
    >
      <g class="sleep-zzz">
        <text
          x="88"
          y="34"
          font-size="10"
          font-weight="bold"
          fill="#a1a1aa"
          font-family="monospace"
        >
          z
        </text>
        <text
          x="96"
          y="26"
          font-size="13"
          font-weight="bold"
          fill="#a1a1aa"
          font-family="monospace"
        >
          Z
        </text>
      </g>
      <g class="sleep-body">
        <circle cx="55" cy="62" r="24" fill="#f4a23b" />
        <path
          d="M40 46 q3 6 0 10 M52 42 q3 7 0 12 M64 46 q3 6 0 10"
          stroke="#d98324"
          stroke-width="3"
          stroke-linecap="round"
          fill="none"
          opacity="0.6"
        />
      </g>
      <g class="sleep-tail">
        <path
          d="M78 70 Q82 84 62 85 Q42 87 34 79"
          fill="none"
          stroke="#d98324"
          stroke-width="7"
          stroke-linecap="round"
        />
      </g>
      <g class="sleep-head">
        <polygon points="62,42 63,32 70,39" fill="#f4a23b" />
        <polygon points="74,39 81,32 82,42" fill="#f4a23b" />
        <circle cx="72" cy="52" r="13" fill="#f4a23b" />
        <ellipse cx="72" cy="57" rx="5.5" ry="4" fill="#ffd9a0" />
        <path
          d="M65 51 h5 M75 51 h5"
          stroke="#2b2b2e"
          stroke-width="1.8"
          stroke-linecap="round"
        />
        <path d="M70.5 55 h3 l-1.5 2 z" fill="#e0607a" />
        <path
          d="M69 58 q3 2 6 0"
          stroke="#2b2b2e"
          stroke-width="1.3"
          fill="none"
          stroke-linecap="round"
        />
      </g>
    </svg>
    <svg v-else viewBox="0 0 150 90" class="h-full w-full pose-in" aria-hidden="true">
      <g transform="translate(74, 80) rotate(-2.2)">
        <g class="play-string">
          <line
            x1="0"
            y1="0"
            x2="51"
            y2="0"
            stroke="#d98324"
            stroke-width="1.5"
            stroke-linecap="round"
          />
        </g>
      </g>
      <g class="play-ball">
        <circle cx="70" cy="78" r="7" fill="#d94f6b" />
        <path
          d="M65 75 q6 3 12 0 M65 81 q6 -3 12 0 M70 72 q3 6 0 12"
          stroke="#a1234b"
          stroke-width="1.3"
          fill="none"
        />
      </g>

      <g class="play-cat">
        <g class="play-tail">
          <path
            d="M25 60 Q6 58 12 36 Q17 26 14 26"
            fill="none"
            stroke="#d98324"
            stroke-width="6"
            stroke-linecap="round"
          />
        </g>

        <g class="play-haunch">
          <circle cx="32" cy="63" r="8" fill="#f4a23b" />
          <circle cx="32" cy="62" r="8" fill="#d98324" opacity="0.25" />
        </g>

        <ellipse cx="52" cy="64" rx="28" ry="13" fill="#f4a23b" />
        <path
          d="M40 54 q3 5 0 9 M52 52 q3 6 0 11"
          stroke="#d98324"
          stroke-width="3"
          stroke-linecap="round"
          fill="none"
          opacity="0.6"
        />

        <rect x="30" y="70" width="7" height="14" rx="3.5" fill="#d98324" />
        <rect x="44" y="72" width="7" height="12" rx="3.5" fill="#f4a23b" />

        <g class="play-leg">
          <rect x="70" y="65" width="7" height="20" rx="3.5" fill="#f4a23b" />
        </g>

        <g class="play-head">
          <polygon points="72,42 74,30 82,38" fill="#f4a23b" />
          <polygon points="85,38 93,30 95,42" fill="#f4a23b" />
          <polygon points="74,40 76,34 80,38" fill="#ef7d8f" />
          <polygon points="87,38 91,34 93,40" fill="#ef7d8f" />
          <circle cx="83" cy="50" r="13" fill="#f4a23b" />
          <ellipse cx="85" cy="55" rx="5.5" ry="4" fill="#ffd9a0" />
          <g class="eyes">
            <circle cx="79" cy="48" r="3.2" fill="#fff" />
            <circle cx="89" cy="48" r="3.2" fill="#fff" />
            <circle cx="80" cy="48.5" r="1.5" fill="#2b2b2e" />
            <circle cx="90" cy="48.5" r="1.5" fill="#2b2b2e" />
          </g>
          <path d="M82.5 52.5 h3 l-1.5 2 z" fill="#e0607a" />
          <path
            d="M81 56 q3 2 6 0"
            stroke="#2b2b2e"
            stroke-width="1.3"
            fill="none"
            stroke-linecap="round"
          />
          <path
            d="M70 52 h-7 M70 55 h-7 M96 52 h7 M96 55 h7"
            stroke="#fff"
            stroke-width="1"
            stroke-linecap="round"
            opacity="0.55"
          />
        </g>
      </g>
    </svg>
  </span>
</template>

<style scoped>
.pet-sprite {
  --ps: var(--pet-speed, 1);
}

.pose-in {
  animation: pose-in 0.18s ease-out;
}

@keyframes pose-in {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.pet-sprite .eyes {
  transform-box: fill-box;
  transform-origin: center;
  animation: blink 5s infinite;
}

.pet-sprite .rig-ear {
  transform-box: fill-box;
  transform-origin: 50% 100%;
  animation: ear-twitch 7s infinite;
}

@keyframes blink {
  0%,
  91%,
  100% {
    transform: scaleY(1);
  }
  94% {
    transform: scaleY(0.12);
  }
  97% {
    transform: scaleY(1);
  }
}

@keyframes ear-twitch {
  0%,
  86%,
  94%,
  100% {
    transform: rotate(0deg);
  }
  89% {
    transform: rotate(-14deg);
  }
  92% {
    transform: rotate(8deg);
  }
}

/* ============ RIG ============ */
.rig-leg {
  transform-box: fill-box;
  transform-origin: 50% 10%;
}
.rig-tail {
  transform-box: fill-box;
  transform-origin: 100% 100%;
}
.rig-body {
  transform-box: fill-box;
  transform-origin: 50% 100%;
}
.rig-head {
  transform-box: fill-box;
  transform-origin: 50% 90%;
}

.stage-walk .rig-leg-fn,
.stage-walk .rig-leg-bf {
  animation: leg-a calc(0.5s / var(--ps)) ease-in-out infinite;
}
.stage-walk .rig-leg-ff,
.stage-walk .rig-leg-bn {
  animation: leg-b calc(0.5s / var(--ps)) ease-in-out infinite;
}
.stage-walk .rig-body {
  animation: body-bob calc(0.25s / var(--ps)) ease-in-out infinite;
}
.stage-walk .rig-head {
  animation: head-bob calc(0.5s / var(--ps)) ease-in-out infinite;
  animation-delay: 0.06s;
}
.stage-walk .rig-tail {
  animation: tail-walk calc(0.5s / var(--ps)) ease-in-out infinite;
}

@keyframes leg-a {
  0%,
  100% {
    transform: rotate(20deg);
  }
  50% {
    transform: rotate(-20deg);
  }
}

@keyframes leg-b {
  0%,
  100% {
    transform: rotate(-20deg);
  }
  50% {
    transform: rotate(20deg);
  }
}

@keyframes body-bob {
  0%,
  100% {
    transform: translateY(0) rotate(0.5deg);
  }
  50% {
    transform: translateY(-2px) rotate(1.2deg);
  }
}

@keyframes head-bob {
  0%,
  100% {
    transform: translateY(0) rotate(0deg);
  }
  50% {
    transform: translateY(-1.5px) rotate(1deg);
  }
}

@keyframes tail-walk {
  0%,
  100% {
    transform: rotate(-8deg);
  }
  50% {
    transform: rotate(12deg);
  }
}

/* ============ ЗУМИС ============ */
.stage-zoomies .rig-leg-fn,
.stage-zoomies .rig-leg-bf {
  animation: leg-a-fast calc(0.22s / var(--ps)) ease-in-out infinite;
}
.stage-zoomies .rig-leg-ff,
.stage-zoomies .rig-leg-bn {
  animation: leg-b-fast calc(0.22s / var(--ps)) ease-in-out infinite;
}
.stage-zoomies .rig-body {
  animation: body-bob calc(0.11s / var(--ps)) ease-in-out infinite;
}
.stage-zoomies .rig-head {
  transform: rotate(4deg) translateX(2px);
}
.stage-zoomies .rig-tail {
  transform: rotate(-18deg);
  animation: tail-fast calc(0.22s / var(--ps)) ease-in-out infinite;
}

@keyframes leg-a-fast {
  0%,
  100% {
    transform: rotate(28deg);
  }
  50% {
    transform: rotate(-28deg);
  }
}

@keyframes leg-b-fast {
  0%,
  100% {
    transform: rotate(-28deg);
  }
  50% {
    transform: rotate(28deg);
  }
}

@keyframes tail-fast {
  0%,
  100% {
    transform: rotate(-14deg);
  }
  50% {
    transform: rotate(-22deg);
  }
}

/* ============ IDLE ============ */
.stage-idle .rig-body {
  animation: breathe calc(2.6s / var(--ps)) ease-in-out infinite;
}
.stage-idle .rig-tail {
  animation: tail-idle calc(1.8s / var(--ps)) ease-in-out infinite;
}
.stage-idle .rig-head {
  animation: head-look calc(7s / var(--ps)) ease-in-out infinite;
}

@keyframes breathe {
  0%,
  100% {
    transform: scaleY(1);
  }
  50% {
    transform: scaleY(1.03);
  }
}

@keyframes tail-idle {
  0%,
  100% {
    transform: rotate(-6deg);
  }
  50% {
    transform: rotate(10deg);
  }
}

@keyframes head-look {
  0%,
  68%,
  92%,
  100% {
    transform: rotate(0deg);
  }
  74%,
  86% {
    transform: rotate(7deg);
  }
}

/* ============ ЗА ШКИРКУ / ПОЛЁТ ============ */
.scruff-sway {
  transform-box: fill-box;
  transform-origin: 50% 0%;
}

.stage-scruff .scruff-sway {
  animation: scruff-sway calc(1.8s / var(--ps)) ease-in-out infinite;
}

.scruff-leg {
  transform-box: fill-box;
  transform-origin: 50% 0%;
}

.stage-scruff .scruff-leg-fl {
  animation: dangle-a calc(1.4s / var(--ps)) ease-in-out infinite;
}
.stage-scruff .scruff-leg-fr {
  animation: dangle-b calc(1.4s / var(--ps)) ease-in-out infinite;
}
.stage-scruff .scruff-leg-bl {
  animation: dangle-b calc(1.7s / var(--ps)) ease-in-out infinite;
}
.stage-scruff .scruff-leg-br {
  animation: dangle-a calc(1.7s / var(--ps)) ease-in-out infinite;
}
.stage-scruff .scruff-tail {
  transform-box: fill-box;
  transform-origin: 100% 0%;
  animation: dangle-a calc(2s / var(--ps)) ease-in-out infinite;
}

@keyframes scruff-sway {
  0%,
  100% {
    transform: rotate(-3deg);
  }
  50% {
    transform: rotate(3deg);
  }
}

@keyframes dangle-a {
  0%,
  100% {
    transform: rotate(6deg);
  }
  50% {
    transform: rotate(16deg);
  }
}

@keyframes dangle-b {
  0%,
  100% {
    transform: rotate(-6deg);
  }
  50% {
    transform: rotate(-16deg);
  }
}

/* в полёте лапы растопырены, не болтаются + легкое трепетание */
.stage-thrown .scruff-sway {
  animation: thrown-flutter calc(0.2s / var(--ps)) ease-in-out infinite alternate;
}
.stage-thrown .scruff-leg-fl {
  transform: rotate(35deg);
}
.stage-thrown .scruff-leg-fr {
  transform: rotate(-35deg);
}
.stage-thrown .scruff-leg-bl {
  transform: rotate(-25deg);
}
.stage-thrown .scruff-leg-br {
  transform: rotate(25deg);
}
.stage-thrown .scruff-tail {
  transform: rotate(-30deg);
}

@keyframes thrown-flutter {
  0% { transform: rotate(-2deg) scale(1, 1); }
  100% { transform: rotate(2deg) scale(1.02, 0.98); }
}

/* Ветер в полёте */
.wind-lines {
  transform-box: fill-box;
  transform-origin: center;
  animation: wind-blow calc(0.4s / var(--ps)) linear infinite;
}

@keyframes wind-blow {
  0% { transform: translateX(15px); opacity: 0; }
  50% { opacity: 0.8; }
  100% { transform: translateX(-25px); opacity: 0; }
}

/* ============ СИДИТ ============ */
.sit-body {
  transform-box: fill-box;
  transform-origin: 50% 100%;
  animation: breathe calc(3s / var(--ps)) ease-in-out infinite;
}
.sit-tail {
  transform-box: fill-box;
  transform-origin: 0% 50%;
  animation: sit-tail-flick calc(5s / var(--ps)) ease-in-out infinite;
}
.sit-head {
  transform-box: fill-box;
  transform-origin: 50% 90%;
  animation: head-look calc(8s / var(--ps)) ease-in-out infinite;
}

@keyframes sit-tail-flick {
  0%,
  76%,
  100% {
    transform: rotate(0deg);
  }
  82% {
    transform: rotate(-9deg);
  }
  88% {
    transform: rotate(4deg);
  }
  94% {
    transform: rotate(0deg);
  }
}

.dizzy-stars {
  transform-box: fill-box;
  transform-origin: center;
  animation: stars-orbit calc(1.2s / var(--ps)) linear infinite;
}

@keyframes stars-orbit {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* ============ СПИТ ============ */
.sleep-body {
  transform-box: fill-box;
  transform-origin: 50% 100%;
  animation: sleep-breathe calc(2.6s / var(--ps)) ease-in-out infinite;
}
.sleep-zzz {
  animation: zzz-float calc(2.2s / var(--ps)) ease-in-out infinite;
}

@keyframes sleep-breathe {
  0%,
  100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.04, 1.06);
  }
}

@keyframes zzz-float {
  0%,
  100% {
    transform: translateY(0);
    opacity: 0.7;
  }
  50% {
    transform: translateY(-4px);
    opacity: 1;
  }
}

/* ============ КЛУБОК ============ */
.play-leg {
  transform-box: fill-box;
  transform-origin: 100% 50%;
  animation: play-swipe calc(2.8s / var(--ps)) ease-in-out infinite;
}
.play-ball {
  transform-box: fill-box;
  transform-origin: center;
  animation: play-ball calc(2.8s / var(--ps)) ease-in-out infinite;
}
.play-string {
  transform-box: fill-box;
  transform-origin: 0% 50%; 
  animation: play-string calc(2.8s / var(--ps)) ease-in-out infinite;
}
.play-cat {
  animation: play-lunge calc(2.8s / var(--ps)) ease-in-out infinite;
}
.play-haunch {
  transform-box: fill-box;
  transform-origin: 50% 100%;
  animation: butt-wiggle calc(2.8s / var(--ps)) ease-in-out infinite;
}
.play-tail {
  transform-box: fill-box;
  transform-origin: 100% 100%;
  animation: tail-fast calc(0.4s / var(--ps)) ease-in-out infinite;
}
.play-head {
  transform-box: fill-box;
  transform-origin: 50% 90%;
  animation: play-head-track calc(2.8s / var(--ps)) ease-in-out infinite;
}

@keyframes play-swipe {
  0%,
  14% {
    transform: rotate(30deg);
  }
  22% {
    transform: rotate(-60deg);
  }
  35%,
  52% {
    transform: rotate(20deg);
  }
  60% {
    transform: rotate(-30deg);
  }
  70%,
  100% {
    transform: rotate(20deg);
  }
}

@keyframes play-ball {
  0%,
  16% {
    transform: translateX(0);
  }
  30% {
    transform: translateX(55px);
  }
  48% {
    transform: translateX(60px);
  }
  60% {
    transform: translateX(-2px);
  }
  68% {
    transform: translateX(4px);
  }
  78%,
  100% {
    transform: translateX(0);
  }
}

@keyframes play-string {
  0%,
  16% {
    transform: scaleX(0.1); /* Клубок рядом с лапой — нитка короткая */
  }
  30%,
  48% {
    transform: scaleX(1); /* Клубок улетел — нитка натянута на всю длину (51px) */
  }
  60%,
  100% {
    transform: scaleX(0.1); /* Клубок вернулся назад */
  }
}

@keyframes play-lunge {
  0%,
  14% {
    transform: translateX(0);
  }
  22% {
    transform: translateX(5px);
  }
  45%,
  52% {
    transform: translateX(0);
  }
  60% {
    transform: translateX(-4px);
  }
  70%,
  100% {
    transform: translateX(0);
  }
}

@keyframes butt-wiggle {
  0%,
  14%,
  100% {
    transform: rotate(0deg);
  }
  3% {
    transform: rotate(-5deg);
  }
  6% {
    transform: rotate(4deg);
  }
  9% {
    transform: rotate(-4deg);
  }
  12% {
    transform: rotate(3deg);
  }
}

@keyframes play-head-track {
  0%,
  16% {
    transform: rotate(0deg);
  }
  30%,
  48% {
    transform: rotate(7deg);
  }
  60% {
    transform: rotate(-4deg);
  }
  75%,
  100% {
    transform: rotate(0deg);
  }
}

@media (prefers-reduced-motion: reduce) {
  .pet-sprite *,
  .pet-sprite {
    animation: none !important;
    transition: none !important;
  }
}
</style>