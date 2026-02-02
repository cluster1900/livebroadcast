import { useEffect, useRef, useState } from 'react'
import videojs from 'video.js'
import 'video.js/dist/video-js.css'
import flvjs from 'flv.js'

interface VideoPlayerProps {
	streamKey: string
	streamURL: string
	poster?: string
	onPlay?: () => void
	onPause?: () => void
	onEnded?: () => void
	onError?: (error: any) => void
}

export function VideoPlayer({
	streamURL,
	poster,
	onPlay,
	onPause,
	onEnded,
	onError
}: VideoPlayerProps) {
	const videoRef = useRef<HTMLVideoElement>(null)
	const flvRef = useRef<flvjs.Player | null>(null)
	const [quality, setQuality] = useState('auto')
	const [qualities] = useState<{ name: string; url: string }[]>([])

	useEffect(() => {
		if (!videoRef.current || !streamURL) return

		const video = videoRef.current
		const isHLS = streamURL.includes('.m3u8')

		if (!isHLS && flvjs.isSupported()) {
			const flv = flvjs.createPlayer({
				type: 'flv',
				url: streamURL
			})

			flv.attachMediaElement(video)
			flv.load()
			flv.play()

			flv.on(flvjs.Events.ERROR, (error: any) => {
				console.error('FLV Error:', error)
				onError?.(error)
			})

			flvRef.current = flv
		} else {
			const player = videojs(videoRef.current, {
				autoplay: true,
				controls: true,
				responsive: true,
				fluid: true,
				poster: poster,
				sources: [
					{
						src: streamURL,
						type: 'application/x-mpegURL'
					}
				]
			})

			player.on('play', () => onPlay?.())
			player.on('pause', () => onPause?.())
			player.on('ended', () => onEnded?.())
			player.on('error', (e: any) => onError?.(e))

			flvRef.current = player as any
		}

		return () => {
			if (flvRef.current) {
				flvRef.current.destroy()
				flvRef.current = null
			}
		}
	}, [streamURL, poster])

	useEffect(() => {
		if (qualities.length > 0 && videoRef.current) {
			// Quality selector logic
		}
	}, [qualities])

	const handleQualityChange = (newQuality: string) => {
		setQuality(newQuality)
	}

	return (
		<div className="video-player-container" style={{ position: 'relative', width: '100%' }}>
			<video
				ref={videoRef}
				className="video-js vjs-big-play-centered"
				controls
				playsInline
				style={{ width: '100%', height: '100%' }}
			/>
			{qualities.length > 0 && (
				<div className="quality-selector" style={{ position: 'absolute', top: 10, right: 10, zIndex: 100 }}>
					<select value={quality} onChange={(e) => handleQualityChange(e.target.value)}>
						<option value="auto">自动</option>
						{qualities.map((q) => (
							<option key={q.name} value={q.name}>{q.name}</option>
						))}
					</select>
				</div>
			)}
		</div>
	)
}

export function FLVPlayer({
	src,
	poster,
	onPlay,
	onPause,
	onEnded,
	onError
}: {
	src: string
	poster?: string
	onPlay?: () => void
	onPause?: () => void
	onEnded?: () => void
	onError?: (error: any) => void
}) {
	const videoRef = useRef<HTMLVideoElement>(null)
	const flvRef = useRef<flvjs.Player | null>(null)

	useEffect(() => {
		if (!videoRef.current || !src) return

		const video = videoRef.current
		const isHLS = src.includes('.m3u8')

		if (!isHLS && flvjs.isSupported()) {
			const flv = flvjs.createPlayer({
				type: 'flv',
				url: src
			})

			flv.attachMediaElement(video)
			flv.load()
			flv.play()

			flv.on(flvjs.Events.ERROR, (error: any) => {
				console.error('FLV Error:', error)
				onError?.(error)
			})

			flvRef.current = flv
		} else {
			const player = videojs(videoRef.current, {
				autoplay: true,
				controls: true,
				responsive: true,
				fluid: true,
				poster: poster,
				sources: [
					{
						src: src,
						type: 'application/x-mpegURL'
					}
				]
			})

			player.on('play', () => onPlay?.())
			player.on('pause', () => onPause?.())
			player.on('ended', () => onEnded?.())
			player.on('error', (e: any) => onError?.(e))

			flvRef.current = player as any
		}

		return () => {
			if (flvRef.current) {
				flvRef.current.destroy()
				flvRef.current = null
			}
		}
	}, [src, poster])

	return (
		<video
			ref={videoRef}
			className="video-js vjs-big-play-centered"
			controls
			playsInline
			style={{ width: '100%', height: '100%' }}
		/>
	)
}

export function HLSPlayer({
	src,
	poster,
	onPlay,
	onPause,
	onEnded,
	onError
}: {
	src: string
	poster?: string
	onPlay?: () => void
	onPause?: () => void
	onEnded?: () => void
	onError?: (error: any) => void
}) {
	const videoRef = useRef<HTMLVideoElement>(null)
	const playerRef = useRef<any>(null)

	useEffect(() => {
		if (!videoRef.current || !src) return

		const player = videojs(videoRef.current, {
			autoplay: true,
			controls: true,
			responsive: true,
			fluid: true,
			poster: poster,
			sources: [
				{
					src: src,
					type: 'application/x-mpegURL'
				}
			]
		})

		playerRef.current = player

		player.on('play', () => onPlay?.())
		player.on('pause', () => onPause?.())
		player.on('ended', () => onEnded?.())
		player.on('error', (e: any) => onError?.(e))

		return () => {
			if (player) {
				player.dispose()
			}
		}
	}, [src, poster])

	return (
		<video
			ref={videoRef}
			className="video-js vjs-big-play-centered"
			controls
			playsInline
			style={{ width: '100%', height: '100%' }}
		/>
	)
}
