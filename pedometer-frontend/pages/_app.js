import "@duik/it/dist/styles.css"; // UI Library: dashboard-ui.com
import "bootstrap/dist/css/bootstrap.css";
import Router from "next/router";
import NProgress from "nprogress"; //nprogress module
import "nprogress/nprogress.css"; //styles of nprogress
import { SWRConfig } from "swr";
import "../styles/styles.css";
import { fetcher } from "../utils/fetcher";

//Binding events.
Router.events.on("routeChangeStart", () => NProgress.start()); Router.events.on("routeChangeComplete", () => NProgress.done()); Router.events.on("routeChangeError", () => NProgress.done());

function MyApp({ Component, pageProps }) {
	return <SWRConfig value={{ fetcher, refreshInterval: 3000 }}>
		<Component {...pageProps} />
	</SWRConfig>;
}

export default MyApp;
