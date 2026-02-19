import React, { useState, useEffect } from 'react';
import { consentService } from '../services/consentService';
import { Cookie, X, ShieldCheck } from 'lucide-react';
import './CookieConsent.css';

const CookieConsent: React.FC = () => {
    const [isVisible, setIsVisible] = useState(false);
    const [isExpanded, setIsExpanded] = useState(false);
    const [preferences, setPreferences] = useState({
        analytics_cookies: false,
        marketing_cookies: false,
        functional_cookies: false,
    });

    useEffect(() => {
        const checkConsent = async () => {
            const consent = await consentService.getConsent();
            if (!consent) {
                setIsVisible(true);
            } else {
                setPreferences({
                    analytics_cookies: consent.analytics_cookies,
                    marketing_cookies: consent.marketing_cookies,
                    functional_cookies: consent.functional_cookies,
                });
            }
        };
        checkConsent();
    }, []);

    const handleAcceptAll = async () => {
        const allAccepted = {
            analytics_cookies: true,
            marketing_cookies: true,
            functional_cookies: true,
        };
        await consentService.updateConsent(allAccepted);
        setPreferences(allAccepted);
        setIsVisible(false);
    };

    const handleSavePreferences = async () => {
        await consentService.updateConsent(preferences);
        setIsVisible(false);
    };

    const handleDeclineAll = async () => {
        const allDeclined = {
            analytics_cookies: false,
            marketing_cookies: false,
            functional_cookies: false,
        };
        await consentService.updateConsent(allDeclined);
        setPreferences(allDeclined);
        setIsVisible(false);
    };

    if (!isVisible) return null;

    return (
        <div className="cookie-consent-container">
            <div className="cookie-consent-card">
                <div className="cookie-consent-header">
                    <div className="cookie-consent-title-row">
                        <div className="cookie-icon-wrapper">
                            <Cookie size={24} />
                        </div>
                        <div>
                            <h3>Cookie Preferences</h3>
                            <p>We use cookies to enhance your experience.</p>
                        </div>
                    </div>
                    <button
                        onClick={() => setIsVisible(false)}
                        className="close-btn"
                        aria-label="Close"
                    >
                        <X size={18} />
                    </button>
                </div>

                {!isExpanded ? (
                    <div className="btn-group" style={{ flexDirection: 'column' }}>
                        <div className="btn-group">
                            <button
                                onClick={handleAcceptAll}
                                className="accept-btn"
                            >
                                Accept All
                            </button>
                            <button
                                onClick={handleDeclineAll}
                                className="decline-btn"
                            >
                                Necessary Only
                            </button>
                        </div>
                        <button
                            onClick={() => setIsExpanded(true)}
                            className="customize-link"
                        >
                            Customize Preferences
                        </button>
                    </div>
                ) : (
                    <div className="preference-list">
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
                            <PreferenceItem
                                title="Necessary"
                                desc="Required for the site to function properly."
                                checked={true}
                                disabled={true}
                                onChange={() => { }}
                            />
                            <PreferenceItem
                                title="Analytics"
                                desc="Help us understand how you use our site."
                                checked={preferences.analytics_cookies}
                                onChange={(v) => setPreferences(p => ({ ...p, analytics_cookies: v }))}
                            />
                            <PreferenceItem
                                title="Functional"
                                desc="Personalization features like theme and language."
                                checked={preferences.functional_cookies}
                                onChange={(v) => setPreferences(p => ({ ...p, functional_cookies: v }))}
                            />
                            <PreferenceItem
                                title="Marketing"
                                desc="Used to deliver relevant advertisements."
                                checked={preferences.marketing_cookies}
                                onChange={(v) => setPreferences(p => ({ ...p, marketing_cookies: v }))}
                            />
                        </div>
                        <div className="btn-group" style={{ paddingTop: '0.5rem' }}>
                            <button
                                onClick={handleSavePreferences}
                                className="accept-btn"
                            >
                                Save Preferences
                            </button>
                            <button
                                onClick={() => setIsExpanded(false)}
                                className="decline-btn"
                            >
                                Back
                            </button>
                        </div>
                    </div>
                )}

                <div className="footer-note">
                    <ShieldCheck size={12} style={{ color: 'var(--success)' }} />
                    <span>Your privacy is our priority. Read our <a href="#" style={{ textDecoration: 'underline' }}>Cookie Policy</a>.</span>
                </div>
            </div>
        </div>
    );
};

interface PreferenceItemProps {
    title: string;
    desc: string;
    checked: boolean;
    disabled?: boolean;
    onChange: (checked: boolean) => void;
}

const PreferenceItem: React.FC<PreferenceItemProps> = ({ title, desc, checked, disabled, onChange }) => (
    <div className="preference-item">
        <div className="preference-info">
            <h4>{title}</h4>
            <p>{desc}</p>
        </div>
        <label className="switch">
            <input
                type="checkbox"
                checked={checked}
                disabled={disabled}
                onChange={(e) => onChange(e.target.checked)}
            />
            <span className="slider"></span>
        </label>
    </div>
);

export default CookieConsent;
