export interface DanmuMessage {
	type: 'danmu'
	timestamp: number
	data: {
		id: string
		user_id: string
		nickname: string
		level: number
		avatar: string
		content: string
		color: string
		badges: string[]
	}
}

export interface GiftMessage {
	type: 'gift'
	timestamp: number
	data: {
		sender: {
			id: string
			nickname: string
			level: number
		}
		gift: {
			id: number
			name: string
			icon: string
			animation: string
		}
		count: number
		combo: number
		total_value: number
	}
}

export interface OnlineCountMessage {
	type: 'online_count'
	timestamp: number
	data: {
		count: number
	}
}

export interface StreamStatusMessage {
	type: 'stream_status'
	timestamp: number
	data: {
		status: 'live' | 'ended' | 'banned'
		reason?: string
	}
}

export interface NotificationMessage {
	type: 'notification'
	timestamp: number
	data: {
		title: string
		content: string
	}
}

export type MessageType = DanmuMessage | GiftMessage | OnlineCountMessage | StreamStatusMessage | NotificationMessage
