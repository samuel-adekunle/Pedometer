import { NavLink, NavPanel, NavSection } from "@duik/it";
import { useRouter } from "next/router";
import { BiHealth } from "react-icons/bi";


function SideNavPanel({ active_item }) {
	const router = useRouter();

	return <div className="side-panel-wrapper">
		<NavPanel className="side-panel">
			<div className="logo-wrapper">
				<span className="logo-icon-wrapper"><BiHealth className="logo-icon" /></span>
				<span className="logo-text">Pedometer</span>
			</div>
			<NavSection>
				<NavLink
					className={active_item == "index" && "active"}
					onClick={() => router.push("/")}
				>
					Dashboard
        </NavLink>
				<NavLink
					className={active_item == "steps" && "active"}
					onClick={() => router.push("/steps")}
				>
					User Stats
        </NavLink>
				<NavLink
					className={active_item == "challenges" && "active"}
					onClick={() => router.push("/challenges")}
				>
					Challenges
        </NavLink>
			</NavSection>
		</NavPanel>
	</div>;
}

export default SideNavPanel;
