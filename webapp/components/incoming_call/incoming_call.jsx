/*
Copyright 2018 Blindside Networks

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import React from 'react';
import PropTypes from 'prop-types';
import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';
const {Modal} = window.ReactBootstrap;
import {Client4} from 'mattermost-redux/client';

export default class IncomingCallPopup extends React.Component {
	static propTypes = {
		show: PropTypes.bool.isRequired,
		theme: PropTypes.object.isRequired,
		siteURL: PropTypes.string.isRequired,
		currentTeam: PropTypes.object.isRequired,
		incomingCall: {
			meetingId: PropTypes.string,
			fromUserID: PropTypes.string,
		},
		actions: PropTypes.shape({
			dismissIncomingCall: PropTypes.func.isRequired,
			getJoinURL: PropTypes.func.isRequired,
		}).isRequired,
	}
	
	constructor(props) {
		super(props);
		this.state = {
			fromUser: null,
		};
	}

	handleClose = () => {
		this.props.actions.dismissIncomingCall();
	};
	
	async componentDidMount() {
		
	}
	
	async componentDidUpdate() {
		if (this.props.show && !this.state.fromUser) {
			const fromUser = await Client4.getUser(this.props.incomingCall.fromUserID);
			this.setState({
				fromUser,
			});	
		}
	}

	getJoinURL = async () => {
		let myurl;
		const userAgent = navigator.userAgent.toLowerCase();
		let myvar;
		//for electron apps
		if (userAgent.indexOf(' electron/') > -1) {
			myurl = await this.props.actions.getJoinURL(this.state.channelId, this.props.incomingCall.meetingId, '');
			myvar = await myurl.data.joinurl.joinURL;
			window.open(myvar);
		} else { //for webapps to circumvent popup blockers
			let newtab = await window.open('about:blank');
			try {
				myurl = await this.props.actions.getJoinURL(this.state.channelId, this.props.incomingCall.meetingId, '');
				myvar = await myurl.data.joinurl.joinURL;
				newtab.location = myvar;
				newtab.focus();
			} catch (e) {
				newtab.close();
			}
		}
	};

	render() {
		const show = this.props.show && this.state.fromUser;
		
		if (!show) {
			return null;
		}

		const style = getStyle(this.props.theme);

		return (
			<Modal show={show} onHide={this.handleClose}>
				<Modal.Header closeButton={true} style={style.header}/>
				<Modal.Body style={style.body}>
					<div>
						<div>
							<img
								src={`${this.props.siteURL}/api/v4/users/${this.props.incomingCall.fromUserID}/image`}
								className="img-responsive img-circle center-block "
							/>
						</div>
						<div style={style.bodyText}>
							<span>
								{'BigBlueButton meeting request from '}
								<a href={`${this.props.siteURL}/${this.props.currentTeam.name}/messages/@${this.state.fromUser.username}`}>{`@${this.state.fromUser.username}`}</a>
							</span>
						</div>
					</div>
				</Modal.Body>
				<Modal.Footer>
					<button type='button' className='btn btn-default' onClick={this.handleClose}>
						{'Close'}
					</button>
					
					<button type='button' className='btn btn-primary pull-left' onClick={this.getJoinURL}>
						{'Join Meeting'}
					</button>
				</Modal.Footer>
			</Modal>
		);
	}
}

const getStyle = makeStyleFromTheme(() => {
	return {
		header: {
			background: 'transparent',
			color: '#0059A5',
			borderStyle: 'none',
			height: '10px',
			minHeight: '28px'
		},
		body: {
			padding: '0px 0px 10px 0px'
		},
		bodyText: {
			textAlign: 'center',
			margin: '20px 0 0 0',
			fontSize: '17px',
			lineHeight: '19px'
		},
	};
});
