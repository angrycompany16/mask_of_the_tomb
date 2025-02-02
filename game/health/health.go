package health

type HealthComponent struct {
	currentHealth float64
	maxHealth     float64
	// Make some way to register death/damage effects
}

func (h *HealthComponent) TakeDamage(damage float64) {
	h.currentHealth -= damage
}

func NewHealthComponent(maxHealth float64) *HealthComponent {
	return &HealthComponent{maxHealth, maxHealth}
}
