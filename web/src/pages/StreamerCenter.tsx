import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { message, Card, Button, Table, Statistic, Row, Col, Modal, Form, Input, Select, Tabs, Tag, Space, Spin } from 'antd'

interface StreamerInfo {
  user_id: string
  stream_key: string
  rtmp_url: string
  status: string
  is_verified: boolean
  total_revenue: number
  follower_count: number
  total_live_duration: number
}

interface LiveRoomInfo {
  id: string
  title: string
  category: string
  cover_url: string
  status: string
  peak_online: number
  total_views: number
  start_at: string
}

interface GiftTransaction {
  id: number
  sender_id: string
  sender_name: string
  gift_name: string
  gift_count: number
  coin_amount: number
  created_at: string
}

function StreamerCenter() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(true)
  const [streamerInfo, setStreamerInfo] = useState<StreamerInfo | null>(null)
  const [liveRoom, setLiveRoom] = useState<LiveRoomInfo | null>(null)
  const [transactions, setTransactions] = useState<GiftTransaction[]>([])
  const [createModalOpen, setCreateModalOpen] = useState(false)
  const [updateModalOpen, setUpdateModalOpen] = useState(false)
  const [activeTab, setActiveTab] = useState('overview')
  const [createForm] = Form.useForm()
  const [updateForm] = Form.useForm()

  const userId = localStorage.getItem('user_id')
  const accessToken = localStorage.getItem('access_token')

  useEffect(() => {
    if (!accessToken) {
      message.warning('请先登录')
      navigate('/login')
      return
    }
    fetchStreamerInfo()
    fetchLiveRoom()
    fetchTransactions()
  }, [accessToken])

  const fetchStreamerInfo = async () => {
    try {
      const response = await axios.get('/api/v1/streamers/me', {
        headers: { Authorization: `Bearer ${accessToken}` }
      })
      if (response.data.code === 0) {
        setStreamerInfo(response.data.data)
      }
    } catch (error) {
      message.error('获取主播信息失败')
    }
  }

  const fetchLiveRoom = async () => {
    try {
      const response = await axios.get('/api/v1/live/rooms', {
        params: { status: 'live' }
      })
      if (response.data.code === 0) {
        const rooms = response.data.data
        const myRoom = rooms.find((r: any) => r.streamer_id === userId)
        setLiveRoom(myRoom || null)
      }
    } catch (error) {
      console.error('获取直播间失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchTransactions = async () => {
    try {
      const response = await axios.get('/api/v1/wallet/transactions', {
        headers: { Authorization: `Bearer ${accessToken}` }
      })
      if (response.data.code === 0) {
        setTransactions(response.data.data)
      }
    } catch (error) {
      console.error('获取交易记录失败')
    }
  }

  const handleCreateRoom = async (values: any) => {
    try {
      const response = await axios.post('/api/v1/rooms', values, {
        headers: { Authorization: `Bearer ${accessToken}` }
      })
      if (response.data.code === 0) {
        message.success('创建直播间成功！')
        setCreateModalOpen(false)
        createForm.resetFields()
        fetchLiveRoom()
      } else {
        message.error(response.data.message || '创建失败')
      }
    } catch (error) {
      message.error('创建失败')
    }
  }

  const handleUpdateRoom = async (values: any) => {
    if (!liveRoom) return
    try {
      const response = await axios.put(`/api/v1/rooms/${liveRoom.id}`, values, {
        headers: { Authorization: `Bearer ${accessToken}` }
      })
      if (response.data.code === 0) {
        message.success('更新直播间成功！')
        setUpdateModalOpen(false)
        fetchLiveRoom()
      } else {
        message.error(response.data.message || '更新失败')
      }
    } catch (error) {
      message.error('更新失败')
    }
  }

  const handleEndRoom = async () => {
    if (!liveRoom) return
    Modal.confirm({
      title: '确认结束直播',
      content: '结束直播后观众将无法观看，确定要结束吗？',
      onOk: async () => {
        try {
          const response = await axios.post(`/api/v1/rooms/${liveRoom.id}/end`, {}, {
            headers: { Authorization: `Bearer ${accessToken}` }
          })
          if (response.data.code === 0) {
            message.success('已结束直播')
            fetchLiveRoom()
          } else {
            message.error(response.data.message || '结束失败')
          }
        } catch (error) {
          message.error('结束失败')
        }
      }
    })
  }

  const handleCopyStreamKey = () => {
    if (streamerInfo?.stream_key) {
      navigator.clipboard.writeText(streamerInfo.stream_key)
      message.success('推流密钥已复制到剪贴板')
    }
  }

  const handleCopyRTMP = () => {
    if (streamerInfo?.rtmp_url) {
      navigator.clipboard.writeText(streamerInfo.rtmp_url)
      message.success('推流地址已复制到剪贴板')
    }
  }

  const formatDuration = (seconds: number) => {
    const hours = Math.floor(seconds / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    return `${hours}小时${minutes}分钟`
  }

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        <Spin size="large" tip="加载中..." />
      </div>
    )
  }

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto', padding: '20px' }}>
      <Card
        title={
          <Space>
            <span>🎮 主播中心</span>
            {streamerInfo?.is_verified && <Tag color="gold">已认证</Tag>}
            {streamerInfo?.status === 'live' && <Tag color="green">直播中</Tag>}
          </Space>
        }
        extra={
          <Button type="primary" onClick={() => navigate('/')}>
            返回首页
          </Button>
        }
      >
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={[
            {
              key: 'overview',
              label: '数据概览',
              children: (
                <Row gutter={16}>
                  <Col span={6}>
                    <Card>
                      <Statistic
                        title="今日收入"
                        value={streamerInfo?.total_revenue || 0}
                        prefix="💰"
                        suffix="虎牙币"
                      />
                    </Card>
                  </Col>
                  <Col span={6}>
                    <Card>
                      <Statistic
                        title="粉丝数量"
                        value={streamerInfo?.follower_count || 0}
                        prefix="👥"
                      />
                    </Card>
                  </Col>
                  <Col span={6}>
                    <Card>
                      <Statistic
                        title="直播时长"
                        value={formatDuration(streamerInfo?.total_live_duration || 0)}
                        prefix="⏱️"
                      />
                    </Card>
                  </Col>
                  <Col span={6}>
                    <Card>
                      <Statistic
                        title="当前状态"
                        value={streamerInfo?.status === 'live' ? '直播中' : '离线'}
                        valueStyle={{ color: streamerInfo?.status === 'live' ? '#52c41a' : '#999' }}
                        prefix={streamerInfo?.status === 'live' ? '🟢' : '🔴'}
                      />
                    </Card>
                  </Col>
                </Row>
              )
            },
            {
              key: 'stream',
              label: '直播管理',
              children: (
                <div>
                  {liveRoom ? (
                    <Card title="当前直播间" style={{ marginBottom: 16 }}>
                      <Row gutter={16}>
                        <Col span={16}>
                          <h3>{liveRoom.title}</h3>
                          <p>分类: {liveRoom.category}</p>
                          <p>观看人数峰值: {liveRoom.peak_online}</p>
                          <p>总观看: {liveRoom.total_views}</p>
                          <Space>
                            <Button type="primary" onClick={() => {
                              updateForm.setFieldsValue({ title: liveRoom.title, category: liveRoom.category })
                              setUpdateModalOpen(true)
                            }}>
                              修改直播间
                            </Button>
                            <Button danger onClick={handleEndRoom}>
                              结束直播
                            </Button>
                          </Space>
                        </Col>
                        <Col span={8}>
                          <Button type="link" onClick={() => navigate(`/live/${liveRoom.id}`)}>
                            进入直播间 →
                          </Button>
                        </Col>
                      </Row>
                    </Card>
                  ) : (
                    <Card style={{ textAlign: 'center', padding: 40 }}>
                      <h3>暂无直播中</h3>
                      <Button type="primary" size="large" onClick={() => setCreateModalOpen(true)}>
                        开启直播
                      </Button>
                    </Card>
                  )}

                  <Card title="推流信息" style={{ marginTop: 16 }}>
                    <p>
                      <strong>推流地址:</strong> {streamerInfo?.rtmp_url}
                      <Button type="link" size="small" onClick={handleCopyRTMP}>复制</Button>
                    </p>
                    <p>
                      <strong>推流密钥:</strong> {streamerInfo?.stream_key?.slice(0, 8)}****
                      <Button type="link" size="small" onClick={handleCopyStreamKey}>复制</Button>
                    </p>
                    <p style={{ color: '#999', fontSize: 12 }}>
                      使用 OBS 或其他推流软件，设置推流地址和密钥即可开始直播
                    </p>
                  </Card>
                </div>
              )
            },
            {
              key: 'revenue',
              label: '收益记录',
              children: (
                <Table
                  dataSource={transactions}
                  rowKey="id"
                  columns={[
                    { title: '时间', dataIndex: 'created_at', render: (t) => new Date(t).toLocaleString() },
                    { title: '描述', dataIndex: 'description' },
                    { title: '金额', dataIndex: 'amount', render: (v) => <span style={{ color: v > 0 ? 'green' : 'red' }}>{v > 0 ? '+' : ''}{v}</span> },
                    { title: '余额', dataIndex: 'balance_after' },
                  ]}
                  pagination={{ pageSize: 10 }}
                />
              )
            }
          ]}
        />
      </Card>

      {/* 创建直播间弹窗 */}
      <Modal
        title="开启直播"
        open={createModalOpen}
        onCancel={() => setCreateModalOpen(false)}
        footer={null}
      >
        <Form form={createForm} layout="vertical" onFinish={handleCreateRoom}>
          <Form.Item
            name="title"
            label="直播间标题"
            rules={[{ required: true, message: '请输入直播间标题', min: 2, max: 200 }]}
          >
            <Input placeholder="给你的直播间起个名字" />
          </Form.Item>
          <Form.Item
            name="category"
            label="直播分类"
            rules={[{ required: true, message: '请选择直播分类' }]}
          >
            <Select
              placeholder="选择直播分类"
              options={[
                { value: '娱乐', label: '娱乐' },
                { value: '游戏', label: '游戏' },
                { value: '音乐', label: '音乐' },
                { value: '舞蹈', label: '舞蹈' },
                { value: '户外', label: '户外' },
                { value: '科技', label: '科技' },
                { value: '体育', label: '体育' },
                { value: '综合', label: '综合' },
              ]}
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              开始直播
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      {/* 修改直播间弹窗 */}
      <Modal
        title="修改直播间"
        open={updateModalOpen}
        onCancel={() => setUpdateModalOpen(false)}
        footer={null}
      >
        <Form form={updateForm} layout="vertical" onFinish={handleUpdateRoom}>
          <Form.Item
            name="title"
            label="直播间标题"
            rules={[{ required: true, message: '请输入直播间标题' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="category"
            label="直播分类"
            rules={[{ required: true, message: '请选择直播分类' }]}
          >
            <Select
              options={[
                { value: '娱乐', label: '娱乐' },
                { value: '游戏', label: '游戏' },
                { value: '音乐', label: '音乐' },
                { value: '舞蹈', label: '舞蹈' },
                { value: '户外', label: '户外' },
                { value: '科技', label: '科技' },
                { value: '体育', label: '体育' },
                { value: '综合', label: '综合' },
              ]}
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              保存修改
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default StreamerCenter
