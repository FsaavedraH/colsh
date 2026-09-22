export default function IlustracionMoto() {
  return (
    <svg
      width="100%"
      height="100%"
      viewBox="0 0 300 200"
      fill="none"
      stroke="#ffffff"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <g id="rueda-trasera">
        <circle cx="70" cy="135" r="32" strokeWidth="4" />
        <circle cx="70" cy="135" r="24" strokeWidth="1.5" />
        <circle cx="70" cy="135" r="14" strokeWidth="2" strokeDasharray="3 4" />
        <circle cx="70" cy="135" r="5" strokeWidth="2" />
        <path d="M70 111 L70 121 M70 159 L70 149 M46 135 L56 135 M94 135 L84 135 M53 118 L60 125 M87 152 L80 145 M87 118 L80 125 M53 152 L60 145" strokeWidth="1.5" />
      </g>

      <g id="rueda-delantera">
        <circle cx="230" cy="135" r="32" strokeWidth="4" />
        <circle cx="230" cy="135" r="24" strokeWidth="1.5" />
        <circle cx="230" cy="135" r="14" strokeWidth="2" strokeDasharray="3 4" />
        <circle cx="230" cy="135" r="5" strokeWidth="2" />
        <path d="M230 111 L230 121 M230 159 L230 149 M206 135 L216 135 M254 135 L244 135 M213 118 L220 125 M247 152 L240 145 M247 118 L240 125 M213 152 L220 145" strokeWidth="1.5" />
      </g>

      <path d="M230 135 L190 55 L180 55 L220 135 Z" />
      <path d="M70 135 L110 130" strokeWidth="2.5" />
      <path d="M95 132 L105 100" strokeWidth="3" strokeDasharray="2 3" />
      <path d="M70 135 L50 110" strokeWidth="1.5" />

      <path d="M100 110 C90 90 60 95 40 120" strokeWidth="2.5" />
      <path d="M40 120 L35 120" strokeWidth="3" />
      <path d="M205 115 C215 100 240 100 255 115" strokeWidth="2.5" />

      <g id="motor">
        <path d="M135 115 L120 85 L140 80 L150 110 Z" fill="#1e3a5f" />
        <path d="M125 90 L140 85 M127 95 L142 90 M130 100 L145 95 M132 105 L147 100" />

        <path d="M150 115 L165 85 L145 80 L135 110 Z" fill="#1e3a5f" />
        <path d="M160 90 L145 85 M157 95 L142 90 M155 100 L140 95 M152 105 L137 100" />

        <circle cx="145" cy="105" r="12" fill="#1e3a5f" />
        <circle cx="145" cy="105" r="5" />

        <path d="M120 130 C130 140 160 140 170 130 C175 120 165 110 145 110 C125 110 115 120 120 130 Z" />
      </g>

      <path d="M145 100 C160 125 130 125 100 125 L45 125" strokeWidth="4" />
      <path d="M45 123 L40 125 L45 127 Z" fill="#ffffff" />
      <path d="M160 110 C175 145 140 140 100 140 L50 140" strokeWidth="4" />
      <path d="M50 138 L45 140 L50 142 Z" fill="#ffffff" />

      <path d="M110 95 C120 70 165 70 185 80 C185 85 160 95 110 95 Z" fill="#1e3a5f" strokeWidth="2.5" />
      <path d="M120 78 C140 70 160 72 170 78" strokeWidth="1" strokeDasharray="4 2" />
      <rect x="155" y="71" width="8" height="3" />

      <path d="M75 100 C80 90 100 95 110 95 C110 95 105 105 85 105 C75 105 70 105 75 100 Z" fill="#1e3a5f" strokeWidth="2" />
      <path d="M60 105 C65 98 75 98 75 100 L75 105 Z" fill="#1e3a5f" />

      <path d="M190 70 C205 70 210 75 210 85 C210 95 205 100 190 100 Z" fill="#1e3a5f" strokeWidth="2" />
      <path d="M205 75 C212 80 212 90 205 95" strokeWidth="1.5" />

      <path d="M185 55 Q170 30 150 40" strokeWidth="2.5" />
      <line x1="150" y1="40" x2="140" y2="43" strokeWidth="4" />

      <line x1="165" y1="42" x2="165" y2="25" strokeWidth="1.5" />
      <circle cx="165" cy="22" r="4" />
    </svg>
  );
}