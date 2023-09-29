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
import SendMail from './components/pages/SendMail.js';
import Help from './components/pages/Help.js';
import Activate from './components/pages/Activate.js';

export default [{
	path: "/", component: Start
},{
	path: "/login", component: Login
},{
	path: "/activate/:userid/:code", component: Activate, props: true
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
	path: "/jobs", component: Jobs,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/node_types", component: NodeTypes,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/node_types/:id", component: NodeTypeDetail, props: true,
	meta: { requiredRole: "ADMIN" }
},{
	path: "/sendmail", component: SendMail,
	meta: { requiredRole: "ADMIN" }
}];
