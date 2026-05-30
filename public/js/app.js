import NavBar from './components/NavBar.js';
import ErrorToast from './components/ErrorToast.js';

export default {
	components: {
		"nav-bar": NavBar,
		"error-toast": ErrorToast
	},
	template: /*html*/`
		<div>
			<nav-bar/>
			<error-toast/>
			<div class="container-fluid">
				<br>
				<div class="alert alert-danger" role="alert">
					<p>
						<b>Server availability:</b> Due to price changes and limited availability on the vps-hoster side
						(<a target="_new" href="https://www.hetzner.com/pressroom/standardization-and-price-adjustment-of-our-server-products/">official post from Hetzner</a>)
						we can't guarantee that new nodes can be created anymore.
					</p>
					<p>
						Existing nodes are still working as expected though, if you are unhappy with the situation
						feel free to ask for a refund in the <router-link to="/tickets">Service-ticket</router-link> page.
					</p>
				</div>
				<router-view></router-view>
			</div>
		</div>
	`
};
