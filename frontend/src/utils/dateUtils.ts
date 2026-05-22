export const formatMessageTime = (timestamp: string): string => {
  if (!timestamp) {
    return new Date().toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
      hour12: true,
    });
  }

  const messageDate = new Date(timestamp);
  const now = new Date();
  const diffInSeconds = Math.floor(
    (now.getTime() - messageDate.getTime()) / 1000
  );
  const diffInMinutes = Math.floor(diffInSeconds / 60);
  const diffInHours = Math.floor(diffInMinutes / 60);
  const diffInDays = Math.floor(diffInHours / 24);

  if (diffInSeconds < 60) {
    return "Just now";
  }

  if (diffInMinutes < 60) {
    return `${diffInMinutes} ${diffInMinutes === 1 ? "minute" : "minutes"} ago`;
  }

  if (diffInHours < 24) {
    return `${diffInHours} ${diffInHours === 1 ? "hour" : "hours"} ago`;
  }

  if (diffInDays < 7) {
    return messageDate.toLocaleString("en-US", {
      weekday: "short",
      hour: "2-digit",
      minute: "2-digit",
      hour12: true,
    });
  }

  return messageDate.toLocaleString("en-US", {
    month: "short",
    day: "numeric",
    year:
      messageDate.getFullYear() !== now.getFullYear() ? "numeric" : undefined,
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  });
};

export const getDateLabel = (timestamp: string): string => {
  if (!timestamp) return "";

  const messageDate = new Date(timestamp);
  const today = new Date();
  const yesterday = new Date(today);
  yesterday.setDate(yesterday.getDate() - 1);

  const messageDateMidnight = new Date(
    messageDate.getFullYear(),
    messageDate.getMonth(),
    messageDate.getDate()
  );
  const todayMidnight = new Date(
    today.getFullYear(),
    today.getMonth(),
    today.getDate()
  );
  const yesterdayMidnight = new Date(
    yesterday.getFullYear(),
    yesterday.getMonth(),
    yesterday.getDate()
  );

  if (messageDateMidnight.getTime() === todayMidnight.getTime()) {
    return "Today";
  } else if (messageDateMidnight.getTime() === yesterdayMidnight.getTime()) {
    return "Yesterday";
  } else {
    return messageDate.toLocaleDateString("en-US", {
      month: "long",
      day: "numeric",
      year:
        messageDate.getFullYear() !== today.getFullYear()
          ? "numeric"
          : undefined,
    });
  }
};

export const isDifferentDay = (
  timestamp1: string,
  timestamp2: string
): boolean => {
  if (!timestamp1 || !timestamp2) return false;

  const date1 = new Date(timestamp1);
  const date2 = new Date(timestamp2);

  return (
    date1.getDate() !== date2.getDate() ||
    date1.getMonth() !== date2.getMonth() ||
    date1.getFullYear() !== date2.getFullYear()
  );
};
