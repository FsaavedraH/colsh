const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface ApiOptions extends RequestInit {
  rol?: string;
}

export async function apiFetch<T>(
  endpoint: string,
  options?: ApiOptions
): Promise<T> {
  const { rol, ...fetchOptions } = options || {};

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(fetchOptions.headers as Record<string, string>),
  };

  if (rol) {
    headers["X-User-Role"] = rol;
  }

  const res = await fetch(`${API_BASE}${endpoint}`, {
    ...fetchOptions,
    headers,
  });

  if (!res.ok) {
    const errorBody = await res.json().catch(() => ({}));
    throw new Error(errorBody.error || `Error ${res.status} en ${endpoint}`);
  }

  return res.json();
}