import CardLayout from "../layouts/CardLayout.js";
import ServerLink from "../ServerLink.js";
import BackupState from "../BackupState.js";
import TimestampBadge from "../TimestampBadge.js";

import { get_all, remove } from "../../api/backup.js";

import format_time from "../../util/format_time.js";
import format_size from "../../util/format_size.js";

export default {
	components: {
		"card-layout": CardLayout,
		"server-link": ServerLink,
		"backup-state": BackupState,
		"timestamp-badge": TimestampBadge
	},
	data: function() {
		return {
			breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "object-group", name: "Backups", link: "/backup"
			}],
			backups: [],
			handle: null
		};
	},
	methods: {
		format_time,
		format_size,
		update_backups: async function() {
			const l = await get_all();
			l.sort((a,b) => a.created < b.created);
			this.backups = l;
		},
		remove_backup: async function(backup) {
			await remove(backup);
			this.update_backups();
		}
	},
	mounted: function() {
		this.handle = setInterval(() => this.update_backups(), 2000);
		this.update_backups();
	},
	beforeUnmount: function() {
		clearInterval(this.handle);
	},
	template: /*html*/`
	<card-layout title="Backups" icon="object-group" :breadcrumb="breadcrumb">
		<table class="table table-condensed table-striped">
			<thead>
				<tr>
					<th>Created</th>
					<th>Expires</th>
					<th>Server</th>
					<th>Size</th>
					<th>State</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="backup in backups" :key="backup.id">
					<td>{{format_time(backup.created)}}</td>
					<td>
						<timestamp-badge :timestamp="backup.expires" :show_duration="true"/>
					</td>
					<td>
						<server-link :id="backup.minetest_server_id"/>
					</td>
					<td>
						<span v-if="backup.state == 'COMPLETE'">
							{{format_size(backup.size)}}
						</span>
					</td>
					<td>
						<backup-state :state="backup.state"/>
					</td>
					<td>
						<div class="btn-group">
							<button class="btn btn-sm btn-danger"
								v-on:click="remove_backup(backup)"
								:disabled="backup.state != 'COMPLETE' && backup.state != 'ERROR'">
								<i class="fa fa-trash"></i>
								Delete
							</button>
							<a class="btn btn-sm btn-secondary" :href="'/api/backup/' + backup.id + '/download'">
								<i class="fa fa-download"></i>
								Download
							</a>
							<router-link class="btn btn-sm btn-outline-secondary" :to="'/mtservers/create?restore_from=' + backup.id">
								<i class="fa fa-play"></i>
								Restore
							</router-link>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
	</card-layout>
	`
};
