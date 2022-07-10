export const API_KEY = "3PaVXHvc7k2aKi3KG3BIy7kaKHAXcAPC27H8uJOG";
export const BASE_URL = "https://2d6mruw8md.execute-api.eu-west-2.amazonaws.com/production";
export const fetcher = (url) => fetch(url, { headers: { "x-api-key": API_KEY } }).then(res => res.json());
