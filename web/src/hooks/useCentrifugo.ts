import { useEffect, useState, useCallback } from 'react'
import { message } from 'antd'
import { DanmuMessage, GiftMessage, OnlineCountMessage, StreamStatusMessage } from '../components/CentrifugoTypes'

interface UseCentrifugoOptions {
	roomId: string
	userId: string
	onDanmu?: (danmu: DanmuMessage) => void
	onGift?: (gift: GiftMessage) => void
	onOnlineCount?: (count: OnlineCountMessage) => void
	onStreamStatus?: (status: StreamStatusMessage) => void
	onConnect?: () => void
	onDisconnect?: () => void
}

export function useCentrifugo({
	roomId,
	userId,
	onDanmu,
	onGift,
	onOnlineCount,
	onStreamStatus,
	onConnect,
	onDisconnect
}: UseCentrifugoOptions) {
	const [connected, setConnected] = useState(false)
	const [connecting, setConnecting] = useState(false)
	const [sub, setSub] = useState<any>(null)

	const connect = useCallback(async () => {
		if (connecting || connected) return

		setConnecting(true)

		try {
			const accessToken = localStorage.getItem('access_token')
			if (!accessToken) {
				throw new Error('No access token')
			}

			const Centrifuge = (await import('centrifuge')).default
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			const centrifuge = new (Centrifuge as any)('ws://localhost:8000/connection/websocket', {
				token: accessToken,
				timeout: 5000
			})

			centrifuge.on('connecting', () => {
				console.log('Centrifugo: connecting...')
			})

			centrifuge.on('connected', () => {
				console.log('Centrifugo: connected')
				setConnected(true)
				setConnecting(false)
				onConnect?.()

				const channel = centrifuge.subscribe(roomId, (ctx: any) => {
					const data = ctx.data

					switch (data.type) {
						case 'danmu':
							onDanmu?.(data as DanmuMessage)
							break
						case 'gift':
							onGift?.(data as GiftMessage)
							break
						case 'online_count':
							onOnlineCount?.(data as OnlineCountMessage)
							break
						case 'stream_status':
							onStreamStatus?.(data as StreamStatusMessage)
							break
					}
				})

				setSub(channel)
			})

			centrifuge.on('disconnected', (ctx: any) => {
				console.log('Centrifugo: disconnected', ctx)
				setConnected(false)
				onDisconnect?.()
			})

			centrifuge.connect()
		} catch (error) {
			console.error('Centrifugo connection error:', error)
			message.error('实时连接失败')
			setConnecting(false)
		}
	}, [roomId, userId, connecting, connected, onConnect, onDisconnect, onDanmu, onGift, onOnlineCount, onStreamStatus])

	const disconnect = useCallback(() => {
		if (sub) {
			sub.unsubscribe()
			setSub(null)
		}
		setConnected(false)
	}, [sub])

	const sendDanmu = useCallback(async (content: string, color?: string) => {
		try {
			const response = await fetch('/api/v1/danmu/send', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('access_token')}`
				},
				body: JSON.stringify({
					room_id: roomId,
					content: content,
					color: color || '#ffffff'
				})
			})

			const data = await response.json()

			if (data.code !== 0) {
				message.error(data.message || '发送弹幕失败')
				return false
			}

			return true
		} catch (error) {
			console.error('Send danmu error:', error)
			message.error('发送弹幕失败')
			return false
		}
	}, [roomId])

	const sendGift = useCallback(async (giftId: number, count: number) => {
		try {
			const response = await fetch('/api/v1/gifts/send', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${localStorage.getItem('access_token')}`
				},
				body: JSON.stringify({
					room_id: roomId,
					gift_id: giftId,
					count: count
				})
			})

			const data = await response.json()

			if (data.code !== 0) {
				message.error(data.message || '赠送礼物失败')
				return false
			}

			return true
		} catch (error) {
			console.error('Send gift error:', error)
			message.error('赠送礼物失败')
			return false
		}
	}, [roomId])

	useEffect(() => {
		return () => {
			disconnect()
		}
	}, [])

	return {
		connected,
		connecting,
		connect,
		disconnect,
		sendDanmu,
		sendGift
	}
}
