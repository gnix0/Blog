import './StatusBar.css';

interface StatusBarProps {
  left?: string;
  center?: string;
  right?: string;
}

export default function StatusBar({ left, center, right }: StatusBarProps) {
  return (
    <div className="statusbar">
      <span className="statusbar-left">{left}</span>
      <span className="statusbar-center">{center}</span>
      <span className="statusbar-right">{right}</span>
    </div>
  );
}
