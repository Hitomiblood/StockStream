import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'

import StockTable from '@/components/StockTable.vue'

const stock = {
  id: '1',
  ticker: 'AAPL',
  target_from: '100',
  target_to: '110',
  company: 'Apple',
  action: 'upgrade',
  brokerage: 'Broker',
  rating_from: 'Hold',
  rating_to: 'Buy',
  time: '2026-02-01T00:00:00Z',
  created_at: '2026-02-01T00:00:00Z',
  updated_at: '2026-02-01T00:00:00Z'
}

describe('components/StockTable', () => {
  it('renders loading skeleton rows', () => {
    const wrapper = mount(StockTable, {
      props: {
        items: [],
        loading: true,
        error: null
      }
    })

    expect(wrapper.findAll('.animate-pulse')).toHaveLength(6)
  })

  it('renders error state', () => {
    const wrapper = mount(StockTable, {
      props: {
        items: [],
        loading: false,
        error: 'failed to load'
      }
    })

    expect(wrapper.text()).toContain('failed to load')
  })

  it('renders empty state', () => {
    const wrapper = mount(StockTable, {
      props: {
        items: [],
        loading: false,
        error: null
      }
    })

    expect(wrapper.text()).toContain('No stocks found for the current criteria.')
  })

  it('emits row-click when a row is clicked', async () => {
    const wrapper = mount(StockTable, {
      props: {
        items: [stock],
        loading: false,
        error: null
      }
    })

    await wrapper.find('tbody tr').trigger('click')

    expect(wrapper.emitted('rowClick')).toBeTruthy()
    expect(wrapper.emitted('rowClick')?.[0]).toEqual([stock])
  })
})
