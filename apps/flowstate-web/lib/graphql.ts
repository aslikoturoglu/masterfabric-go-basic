const API_URL = process.env.NEXT_PUBLIC_GRAPHQL_URL || "http://localhost:8080/graphql";

export async function gql<T>(query: string, variables?: Record<string, unknown>, token?: string): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(API_URL, {
    method: "POST",
    headers,
    body: JSON.stringify({ query, variables }),
  });
  const json = await res.json();
  if (json.errors?.length) throw new Error(json.errors[0]?.message || "GraphQL error");
  return json.data as T;
}
