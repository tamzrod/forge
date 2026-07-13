import { useMemo } from 'react';
import type { State } from '../../types';
import styles from './AnalysisPanel.module.css';

type AnalysisTab = 'timeline' | 'events' | 'why';

interface AnalysisPanelProps {
  simulationState: State;
  activeTab: AnalysisTab;
  onTabChange: (tab: AnalysisTab) => void;
}

export function AnalysisPanel({
  simulationState,
  activeTab,
  onTabChange,
}: AnalysisPanelProps) {
  const { clock, sun, weather, grid } = simulationState;

  const tabs: { id: AnalysisTab; label: string }[] = [
    { id: 'timeline', label: 'Timeline' },
    { id: 'events', label: 'Events' },
    { id: 'why', label: 'Why?' },
  ];

  // Generate timeline events
  const timelineEvents = useMemo(() => {
    const events = [];

    // Sunrise/sunset events based on sun state
    if (sun.elevation > 0 && !sun.is_daytime) {
      events.push({
        type: 'sunrise',
        title: 'Sunrise',
        description: 'Sun has risen above horizon. PV generation possible.',
        time: clock.elapsed,
      });
    } else if (sun.elevation <= 0 && sun.is_daytime) {
      events.push({
        type: 'sunset',
        title: 'Sunset',
        description: 'Sun has set below horizon. PV generation ended.',
        time: clock.elapsed,
      });
    }

    // Grid stability events
    if (!grid.is_stable) {
      events.push({
        type: 'alarm',
        title: 'Grid Instability',
        description: 'Grid frequency or voltage outside normal parameters.',
        time: clock.elapsed,
      });
    }

    // Weather events
    if (weather.cloud_cover > 80) {
      events.push({
        type: 'info',
        title: 'Heavy Cloud Cover',
        description: `Cloud cover at ${weather.cloud_cover.toFixed(0)}%. Expect reduced generation.`,
        time: clock.elapsed,
      });
    }

    if (weather.wind_speed > 10) {
      events.push({
        type: 'info',
        title: 'High Wind Advisory',
        description: `Wind speed ${weather.wind_speed.toFixed(1)} m/s. Monitor equipment.`,
        time: clock.elapsed,
      });
    }

    // Add simulation start event
    if (clock.tick_count === 0) {
      events.push({
        type: 'info',
        title: 'Simulation Started',
        description: 'Simulation initialized. Click Run to begin.',
        time: '00:00:00',
      });
    }

    return events.sort((a, b) => b.time.localeCompare(a.time));
  }, [clock.elapsed, clock.tick_count, sun, weather, grid]);

  // Generate event log
  const eventLog = useMemo(() => {
    const events = [];

    // Grid events
    if (grid.voltage_pu < 0.95) {
      events.push({
        type: 'warning',
        message: 'Grid voltage below 95% of nominal',
        source: 'PCC Meter',
        timestamp: clock.elapsed,
      });
    }

    if (!grid.is_stable) {
      events.push({
        type: 'critical',
        message: 'Grid frequency deviation detected',
        source: 'Protection System',
        timestamp: clock.elapsed,
      });
    }

    // Sun events
    if (!sun.is_daytime) {
      events.push({
        type: 'info',
        message: 'Nighttime - no solar generation',
        source: 'Sun Model',
        timestamp: clock.elapsed,
      });
    } else if (sun.irradiance < 100) {
      events.push({
        type: 'info',
        message: 'Low irradiance conditions',
        source: 'Weather Station',
        timestamp: clock.elapsed,
      });
    }

    // Weather events
    if (weather.is_raining) {
      events.push({
        type: 'warning',
        message: 'Rain detected - possible generation reduction',
        source: 'Weather Station',
        timestamp: clock.elapsed,
      });
    }

    if (weather.temperature > 45) {
      events.push({
        type: 'warning',
        message: 'High ambient temperature - monitor equipment',
        source: 'Weather Station',
        timestamp: clock.elapsed,
      });
    }

    // Normal operation event
    if (sun.is_daytime && grid.is_stable && events.length === 0) {
      events.push({
        type: 'success',
        message: 'All systems operating normally',
        source: 'System',
        timestamp: clock.elapsed,
      });
    }

    return events;
  }, [clock.elapsed, sun, weather, grid]);

  // Generate Why explanations
  const whyInsights = useMemo(() => {
    const insights = [];

    // Why is power low?
    if (sun.irradiance < 500) {
      insights.push({
        icon: '☀️',
        title: 'Why is power generation low?',
        content: `The current irradiance is ${sun.irradiance.toFixed(0)} W/m², which is ${(sun.irradiance / 10).toFixed(0)}% of peak conditions. This is the primary factor limiting power output.

Factors affecting irradiance:
• Time of day: ${sun.elevation > 0 ? 'Sun is above horizon' : 'Sun is below horizon'}
• Cloud cover: ${weather.cloud_cover.toFixed(0)}%
• Season: Summer peak generation period`,
        insight: `At 1000 W/m², this plant can generate up to rated capacity. At ${sun.irradiance.toFixed(0)} W/m², expect approximately ${(sun.irradiance / 10).toFixed(0)}% of rated output.`,
      });
    }

    // Why is voltage low?
    if (grid.voltage_pu < 0.98) {
      insights.push({
        icon: '⚡',
        title: 'Why is grid voltage low?',
        content: `The grid voltage is ${grid.voltage.toFixed(0)} V (${(grid.voltage_pu * 100).toFixed(1)}% of nominal).

Possible causes:
• High load on the distribution feeder
• Transformer tap position
• Power factor of local loads
• Grid impedance`,
        insight: `Voltage within ±10% of nominal is acceptable. Below 90% may trigger grid protection.`,
      });
    }

    // Why temperature affects efficiency
    if (weather.temperature > 35) {
      insights.push({
        icon: '🌡️',
        title: 'Why does temperature matter?',
        content: `Current ambient temperature is ${weather.temperature.toFixed(1)}°C.

High temperatures reduce PV panel efficiency:
• Every 1°C above 25°C reduces efficiency by ~0.4%
• Current efficiency loss: ~${((weather.temperature - 25) * 0.4).toFixed(1)}%
• Panel temperature may be ${(weather.temperature + 25).toFixed(0)}°C or higher`,
        insight: `On hot days, generation may be 5-10% lower than expected. Consider cleaning panels to improve cooling.`,
      });
    }

    // Why cloud cover matters
    if (weather.cloud_cover > 30) {
      insights.push({
        icon: '☁️',
        title: 'Why does cloud cover reduce generation?',
        content: `Cloud cover is currently ${weather.cloud_cover.toFixed(0)}%.

Cloud cover directly affects irradiance:
• Clear sky: Up to 1000 W/m²
• 50% clouds: ~500-700 W/m²
• Overcast: 100-300 W/m²
• Heavy clouds: <100 W/m²`,
        insight: `Cloud transients can cause rapid power fluctuations. This is normal and the inverter will handle the changes.`,
      });
    }

    // Always show general info
    if (insights.length === 0) {
      insights.push({
        icon: '📊',
        title: 'Current Generation Analysis',
        content: `Operating conditions:
• Irradiance: ${sun.irradiance.toFixed(0)} W/m²
• Grid: ${grid.voltage.toFixed(0)} V @ ${grid.frequency.toFixed(2)} Hz
• Temperature: ${weather.temperature.toFixed(1)}°C`,
        insight: `All parameters within normal range. Plant is operating optimally.`,
      });
    }

    return insights;
  }, [sun, weather, grid]);

  // Calculate stats
  const stats = useMemo(() => {
    const todayGeneration = sun.irradiance > 0 
      ? (clock.elapsed_ms / 3600000) * (sun.irradiance / 10) * 500 / 1000 
      : 0;
    
    return {
      totalEvents: eventLog.length,
      activeAlarms: eventLog.filter(e => e.type === 'critical' || e.type === 'warning').length,
      todayGeneration: todayGeneration.toFixed(1),
      efficiency: sun.irradiance > 0 ? (sun.irradiance / 10).toFixed(0) : '0',
    };
  }, [clock.elapsed_ms, sun.irradiance, eventLog]);

  // Tab bar shared by all panels
  const tabBar = (
    <div style={{ display: 'flex', gap: '8px', marginBottom: '16px' }}>
      {tabs.map((tab) => (
        <button
          key={tab.id}
          onClick={() => onTabChange(tab.id)}
          style={{
            flex: 1,
            padding: '8px 12px',
            background: activeTab === tab.id ? '#ff6f00' : '#3c3c3c',
            border: 'none',
            borderRadius: '4px',
            color: activeTab === tab.id ? 'white' : '#d4d4d4',
            cursor: 'pointer',
            fontSize: '12px',
            fontWeight: activeTab === tab.id ? 600 : 400,
            transition: 'all 0.15s ease',
          }}
        >
          {tab.label}
        </button>
      ))}
    </div>
  );

  switch (activeTab) {
    case 'timeline':
      return (
        <div className={styles.panel}>
          {tabBar}
          <div className={styles.stats}>
            <div className={styles.stat}>
              <div className={styles.statValue}>{stats.activeAlarms}</div>
              <div className={styles.statLabel}>Active Alarms</div>
            </div>
            <div className={styles.stat}>
              <div className={styles.statValue}>{stats.totalEvents}</div>
              <div className={styles.statLabel}>Total Events</div>
            </div>
          </div>

          <div className={styles.sectionTitle}>Timeline</div>
          <div className={styles.timeline}>
            {timelineEvents.map((event, i) => (
              <div key={i} className={styles.timelineItem}>
                <div className={`${styles.timelineIcon} ${styles[event.type]}`}>
                  {event.type === 'sunrise' && '🌅'}
                  {event.type === 'sunset' && '🌇'}
                  {event.type === 'alarm' && '⚠️'}
                  {event.type === 'info' && 'ℹ️'}
                </div>
                <div className={styles.timelineContent}>
                  <div className={styles.timelineTitle}>{event.title}</div>
                  <div className={styles.timelineDescription}>{event.description}</div>
                  <div className={styles.timelineTime}>{event.time}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      );

    case 'events':
      return (
        <div className={styles.panel}>
          {tabBar}
          <div className={styles.stats}>
            <div className={styles.stat}>
              <div className={styles.statValue}>{stats.activeAlarms}</div>
              <div className={styles.statLabel}>Active Alarms</div>
            </div>
            <div className={styles.stat}>
              <div className={styles.statValue}>{stats.todayGeneration}</div>
              <div className={styles.statLabel}>MWh Today</div>
            </div>
          </div>

          <div className={styles.sectionTitle}>Event Log</div>
          <div className={styles.events}>
            {eventLog.map((event, i) => (
              <div key={i} className={`${styles.event} ${styles[event.type]}`}>
                <div className={styles.eventHeader}>
                  <span className={`${styles.eventType} ${styles[event.type]}`}>
                    {event.type}
                  </span>
                  <span className={styles.eventTimestamp}>{event.timestamp}</span>
                </div>
                <div className={styles.eventMessage}>{event.message}</div>
                <div className={styles.eventSource}>Source: {event.source}</div>
              </div>
            ))}
          </div>
        </div>
      );

    case 'why':
      return (
        <div className={styles.panel}>
          {tabBar}
          <div className={styles.sectionTitle}>Why?</div>
          <div className={styles.why}>
            {whyInsights.map((insight, i) => (
              <div key={i} className={styles.whyCard}>
                <div className={styles.whyCardHeader}>
                  <span className={styles.whyCardIcon}>{insight.icon}</span>
                  <span className={styles.whyCardTitle}>{insight.title}</span>
                </div>
                <div className={styles.whyCardContent}>
                  {insight.content.split('\n\n').map((para, j) => (
                    <p key={j} style={{ marginBottom: para.includes('\n') ? '0' : '8px' }}>
                      {para.split('\n').map((line, k) => (
                        <span key={k}>
                          {line}
                          {k < para.split('\n').length - 1 && <br />}
                        </span>
                      ))}
                    </p>
                  ))}
                </div>
                {insight.insight && (
                  <div className={styles.whyInsight}>
                    <div className={styles.whyInsightTitle}>Key Insight</div>
                    <div className={styles.whyInsightContent}>
                      {insight.insight.split(' ').map((word, j) => (
                        word.includes(':') ? (
                          <span key={j}> <span className={styles.whyHighlight}>{word}</span></span>
                        ) : (
                          <span key={j}> {word}</span>
                        )
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      );

    default:
      return null;
  }
}
