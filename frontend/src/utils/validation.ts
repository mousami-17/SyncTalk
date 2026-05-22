import type { PasswordStrength } from "../types";

export const validateUsername = (
  username: string
): { isValid: boolean; error?: string } => {
  const trimmed = username.trim();

  if (trimmed.length < 3) {
    return {
      isValid: false,
      error: "Username must be at least 3 characters long",
    };
  }

  if (trimmed.length > 30) {
    return { isValid: false, error: "Username must not exceed 30 characters" };
  }

  const usernameRegex = /^[a-zA-Z0-9_-]+$/;
  if (!usernameRegex.test(trimmed)) {
    return {
      isValid: false,
      error:
        "Username can only contain letters, numbers, underscores, and hyphens",
    };
  }

  return { isValid: true };
};

export const calculatePasswordStrength = (
  password: string
): PasswordStrength => {
  if (!password) {
    return {
      score: 0,
      label: "",
      color: "#6b7280",
      requirements: {
        length: false,
        uppercase: false,
        lowercase: false,
        number: false,
        special: false,
      },
    };
  }

  const requirements = {
    length: password.length >= 8,
    uppercase: /[A-Z]/.test(password),
    lowercase: /[a-z]/.test(password),
    number: /[0-9]/.test(password),
    special: /[^A-Za-z0-9]/.test(password),
  };

  let score = 0;
  if (requirements.length) score++;
  if (requirements.uppercase) score++;
  if (requirements.lowercase) score++;
  if (requirements.number) score++;
  if (requirements.special) score++;

  let label: string;
  let color: string;

  if (score <= 2) {
    label = "Weak";
    color = "#ef4444";
  } else if (score === 3) {
    label = "Fair";
    color = "#f97316";
  } else if (score === 4) {
    label = "Good";
    color = "#3b82f6";
  } else {
    label = "Strong";
    color = "#10b981";
  }

  return {
    score,
    label,
    color,
    requirements,
  };
};

export const validatePassword = (
  password: string
): { isValid: boolean; error?: string } => {
  if (password.length < 8) {
    return {
      isValid: false,
      error: "Password must be at least 8 characters long",
    };
  }

  if (password.length > 100) {
    return { isValid: false, error: "Password must not exceed 100 characters" };
  }

  const hasUpper = /[A-Z]/.test(password);
  const hasLower = /[a-z]/.test(password);
  const hasNumber = /[0-9]/.test(password);

  if (!hasUpper || !hasLower || !hasNumber) {
    return {
      isValid: false,
      error:
        "Password must contain at least one uppercase letter, one lowercase letter, and one number",
    };
  }

  return { isValid: true };
};

export const sanitizeMessage = (message: string): string => {
  return message
    .trim()
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;")
    .replace(/\//g, "&#x2F;");
};

export const validateRoomName = (
  roomName: string
): { isValid: boolean; error?: string } => {
  const trimmed = roomName.trim();

  if (!trimmed) {
    return { isValid: false, error: "Room name cannot be empty" };
  }

  if (trimmed.length > 50) {
    return { isValid: false, error: "Room name must not exceed 50 characters" };
  }

  return { isValid: true };
};
