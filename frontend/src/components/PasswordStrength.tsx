import React from "react";
import { Check, X } from "lucide-react";
import { calculatePasswordStrength } from "../utils/validation";

interface PasswordStrengthProps {
  password: string;
}

const PasswordStrength: React.FC<PasswordStrengthProps> = ({ password }) => {
  const strength = calculatePasswordStrength(password);

  if (!password) return null;

  return (
    <div className="space-y-2 mt-2">
      <div className="flex items-center gap-2">
        <div className="flex-1 h-2 bg-gray-700/40 rounded-full overflow-hidden">
          <div
            className="h-full transition-all duration-300 rounded-full"
            style={{
              width: `${(strength.score / 5) * 100}%`,
              backgroundColor: strength.color,
            }}
          />
        </div>
        <span className="text-sm font-medium" style={{ color: strength.color }}>
          {strength.label}
        </span>
      </div>

      <div className="space-y-1 text-xs">
        <div
          className={`flex items-center gap-1 ${
            strength.requirements.length ? "text-gray-300" : "text-white/50"
          }`}
        >
          {strength.requirements.length ? (
            <Check className="w-3 h-3" />
          ) : (
            <X className="w-3 h-3" />
          )}
          <span>At least 8 characters</span>
        </div>
        <div
          className={`flex items-center gap-1 ${
            strength.requirements.uppercase ? "text-gray-300" : "text-white/50"
          }`}
        >
          {strength.requirements.uppercase ? (
            <Check className="w-3 h-3" />
          ) : (
            <X className="w-3 h-3" />
          )}
          <span>One uppercase letter</span>
        </div>
        <div
          className={`flex items-center gap-1 ${
            strength.requirements.lowercase ? "text-gray-300" : "text-white/50"
          }`}
        >
          {strength.requirements.lowercase ? (
            <Check className="w-3 h-3" />
          ) : (
            <X className="w-3 h-3" />
          )}
          <span>One lowercase letter</span>
        </div>
        <div
          className={`flex items-center gap-1 ${
            strength.requirements.number ? "text-gray-300" : "text-white/50"
          }`}
        >
          {strength.requirements.number ? (
            <Check className="w-3 h-3" />
          ) : (
            <X className="w-3 h-3" />
          )}
          <span>One number</span>
        </div>
      </div>
    </div>
  );
};

export default PasswordStrength;
