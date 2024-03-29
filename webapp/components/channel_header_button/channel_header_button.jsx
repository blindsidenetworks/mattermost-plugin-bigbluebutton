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
import {Svgs} from '../../constants';

import PropTypes from 'prop-types';
import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

const {OverlayTrigger, Tooltip} = window.ReactBootstrap;


export default class ChannelHeaderButton extends React.PureComponent {
	static propTypes = {
		channelId: PropTypes.string.isRequired,
		theme: PropTypes.object.isRequired,
		visible: PropTypes.bool.isRequired,
	};

	constructor(props) {
		super(props);

		this.state = {
			showPopover: false
		};

		this.overlayRef = React.createRef();
	}

	render() {
		if (this.props.channelId === '') {
			return <div/>;
		}

		const style = getStyle(this.props.theme);
		return (<div>
			<div>
				<OverlayTrigger
					trigger={['hover']}
					delayShow={400}
					ref={el => this.overlayRef = el}
					placement="bottom"
					overlay={(
						<Tooltip id="bbbChannelHeaderTooltip">
							{'BigBlueButton'}
						</Tooltip>
					)}
				>
					<div
						id="bbbChannelHeaderButton"
						onClick={(e) => {
							this.overlayRef.hide();
							this.setState({
								popoverTarget: e.target,
								showPopover: !this.props.visible
							});
						}}
						style={style.svg}
					>
						<span
							style={style.iconStyle}
							aria-hidden="true"
							dangerouslySetInnerHTML={{
								__html: Svgs.BBB_LOGO_SIMPLIFIED
							}}/>
					</div>
				</OverlayTrigger>
			</div>
		</div>);
	}
}


const getStyle = makeStyleFromTheme((theme) => {
	return {
		iconStyle: {
			position: 'relative',
		},
		svg: {
			width: '20px',
			height: '20px',
		}
	};
});
