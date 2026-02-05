import Start from './components/pages/Start.js';
import Profile from './components/pages/Profile.js';
import Login from './components/pages/Login.js';
import NodeTypes from './components/pages/NodeTypes.js';
import NodeTypeDetail from './components/pages/NodeTypeDetail.js';
import MTServers from './components/pages/MTServers.js';
import MTServerDetail from './components/pages/MTServerDetail.js';
import MTServerCreate from './components/pages/MTServerCreate.js';
import MTServerDelete from './components/pages/MTServerDelete.js';
import UserNodes from './components/pages/UserNodes.js';
import UserNodesDetail from './components/pages/UserNodesDetail.js';
import UserNodeCreate from './components/pages/UserNodeCreate.js';
import UserNodeDelete from './components/pages/UserNodeDelete.js';
import Jobs from './components/pages/Jobs.js';
import Finance from './components/pages/Finance.js';
import FinanceDetail from './components/pages/FinanceDetail.js';
import AuditLogs from './components/pages/AuditLogs.js';
import Help from './components/pages/Help.js';
import Pricing from './components/pages/Pricing.js';
import Register from './components/pages/Register.js';
import Overview from './components/pages/Overview.js';
import BackupList from './components/pages/BackupList.js';
import Users from './components/pages/Users.js';
import UserDetail from './components/pages/UserDetail.js';
import Coupon from './components/pages/Coupon.js';
import ServiceTickets from './components/pages/ticket/ServiceTickets.js';
import NewTicket from './components/pages/ticket/NewTicket.js';
import ServiceTicketDetail from './components/pages/ticket/ServiceTicketDetail.js';

export default [{
	path: "/", component: Start
},{
	path: "/login", component: Login
},{
	path: "/register", component: Register
},{
	path: "/overview", component: Overview
},{
	path: "/pricing", component: Pricing
},{
	path: "/help", component: Help,
	meta: { loggedIn: true }
},{
	path: "/profile", component: Profile,
	meta: { loggedIn: true }
},{
	path: "/audit-logs", component: AuditLogs,
	meta: { loggedIn: true }
},{
	path: "/finance", component: Finance,
	meta: { loggedIn: true }
},{
	path: "/finance/detail/:id", component: FinanceDetail, props: true,
	meta: { loggedIn: true }
},{
	path: "/nodes", component: UserNodes,
	meta: { loggedIn: true }
},{
	path: "/nodes/create", component: UserNodeCreate,
	meta: { loggedIn: true }
},{
	path: "/nodes/:id", component: UserNodesDetail, props: true,
	meta: { loggedIn: true }
},{
	path: "/nodes/:id/delete", component: UserNodeDelete, props: true,
	meta: { loggedIn: true }
},{
	path: "/mtservers", component: MTServers,
	meta: { loggedIn: true }
},{
	path: "/mtservers/create", component: MTServerCreate,
	meta: { loggedIn: true }
},{
	path: "/mtservers/:id", component: MTServerDetail, props: true,
	meta: { loggedIn: true }
},{
	path: "/mtservers/:id/delete", component: MTServerDelete, props: true,
	meta: { loggedIn: true }
},{
	path: "/users", component: Users,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/users/:id", component: UserDetail, props: true,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/jobs", component: Jobs,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/node_types", component: NodeTypes,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/node_types/:id", component: NodeTypeDetail, props: true,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/backup", component: BackupList,
	meta: { loggedIn: true }
},{
	path: "/coupon", component: Coupon,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/tickets", component: ServiceTickets,
	meta: { loggedIn: true }
},{
	path: "/ticket/:id", component: ServiceTicketDetail, props: true,
	meta: { loggedIn: true }
},{
	path: "/tickets/new", component: NewTicket,
	meta: { loggedIn: true }
}];
